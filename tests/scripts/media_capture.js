(() => {
    if (window.__MediaCaptureController) return;
    class MediaCaptureController {
        constructor() {
            this.videoMap = new Map(); // key: videoEl, value: { audioCtx, src, compressor, canvas, ctx, rafId }
        }

        async start(videoEl = null) {
            if (!videoEl) videoEl = document.querySelector('video');
            if (!videoEl) throw new Error('no video element');
            if (this.videoMap.has(videoEl)) return; // 已经在采集

            const video = videoEl;

            // ================== AUDIO ==================
            video.muted = false;
            video.volume = 1.0;
            video.crossOrigin = 'anonymous';

            // 推荐方案配置
            const AUDIO_CONFIG = {
                sampleRate: 24000,       // 24kHz适合语音
                bitDepth: 16,           // 16位整型
                enableVAD: true,        // 语音活动检测
                vadThreshold: 0.001,    // VAD阈值
                silenceFrames: 10,      // 连续静音帧数阈值
                bufferSize: 2048        // 缓冲区大小
            };

            const audioCtx = new AudioContext({ sampleRate: AUDIO_CONFIG.sampleRate });
            await audioCtx.resume();

            try {
                // 定义内联 AudioWorklet 处理器
                const workletCode = `
                class AudioCompressorProcessor extends AudioWorkletProcessor {
                    constructor() {
                        super();
                        this.energyThreshold = ${AUDIO_CONFIG.vadThreshold};
                        this.silenceCount = 0;
                        this.maxSilenceFrames = ${AUDIO_CONFIG.silenceFrames};
                        this.isActive = false;
                    }

                    process(inputs) {
                        const input = inputs[0];
                        if (!input || input.length === 0) return true;
                        
                        const ch0 = input[0];
                        
                        // 语音活动检测
                        const shouldSend = this.vadCheck(ch0);
                        
                        if (shouldSend) {
                            // 32位浮点转16位整型
                            const int16Buffer = this.float32ToInt16(ch0);
                            
                            this.port.postMessage({
                                sampleRate: ${AUDIO_CONFIG.sampleRate},
                                channels: 1,
                                frames: ch0.length,
                                buffer: int16Buffer.buffer,
                                ts: currentTime
                            });
                        }
                        
                        return true;
                    }
                    
                    // 语音活动检测
                    vadCheck(data) {
                        if (!${AUDIO_CONFIG.enableVAD}) return true;
                        
                        let energy = 0;
                        // 只检查前100个样本以提升性能
                        for (let i = 0; i < Math.min(data.length, 100); i++) {
                            energy += Math.abs(data[i]);
                            if (energy > this.energyThreshold) break;
                        }
                        
                        if (energy > this.energyThreshold) {
                            this.silenceCount = 0;
                            this.isActive = true;
                            return true;
                        } else {
                            this.silenceCount++;
                            if (this.silenceCount > this.maxSilenceFrames) {
                                this.isActive = false;
                            }
                            // 静音期间也发送几帧，避免语音截断
                            return this.isActive;
                        }
                    }
                    
                    // 32位浮点转16位整型
                    float32ToInt16(floatArray) {
                        const int16 = new Int16Array(floatArray.length);
                        for (let i = 0; i < floatArray.length; i++) {
                            // 将 [-1, 1] 的浮点数映射到 [-32768, 32767] 的整数
                            const sample = floatArray[i];
                            int16[i] = sample < -1 ? -32768 : 
                                    sample > 1 ? 32767 : 
                                    Math.floor(sample * 32768);
                        }
                        return int16;
                    }
                }
                
                registerProcessor('audio-compressor-processor', AudioCompressorProcessor);
            `;

                // 创建 Blob URL 加载 worklet
                const blob = new Blob([workletCode], { type: 'application/javascript' });
                const blobURL = URL.createObjectURL(blob);
                await audioCtx.audioWorklet.addModule(blobURL);
                URL.revokeObjectURL(blobURL);

                const src = audioCtx.createMediaElementSource(video);
                const compressor = new AudioWorkletNode(audioCtx, 'audio-compressor-processor');

                src.connect(compressor);
                compressor.connect(audioCtx.destination);

                const onAudioMessage = (e) => {
                    if (!this.videoMap.has(video)) return;
                    if (video.paused) return;

                    const audioData = e.data;
                    console.log(audioData.buffer);
                    // 保持原有接口格式，但数据已压缩
                    window?.__onVideoAudio?.({
                        sampleRate: audioData.sampleRate,
                        channels: 1,
                        frames: audioData.frames,
                        buffer: audioData.buffer,  // 已经是16位整型
                        ts: performance.now() / 1000,  // 使用更高精度的时间戳
                        compressed: true,  // 添加标记表明是压缩数据
                        bitDepth: AUDIO_CONFIG.bitDepth
                    });
                };

                compressor.port.onmessage = onAudioMessage.bind(this);

                // ================== VIDEO ==================
                const canvas = document.createElement('canvas');
                const ctx = canvas.getContext('2d');
                let rafId = null;
                let lastCaptureTime = 0;
                const FPS = 1; // 每秒1帧
                const INTERVAL = 1000 / FPS; // 1000ms

                const captureFrame = () => {
                    if (!this.videoMap.has(video)) return;
                    const now = performance.now();
                    if (now - lastCaptureTime >= INTERVAL) {
                        if (video.readyState >= 2 && !video.paused) {
                            const w = video.videoWidth;
                            const h = video.videoHeight;
                            if (w && h) {
                                // 设置目标分辨率（示例：缩小到原尺寸的20%）
                                const scale = 0.2; // 缩放比例，0.2表示20%
                                const targetWidth = Math.floor(w * scale);
                                const targetHeight = Math.floor(h * scale);
                                // 设置canvas为缩小的尺寸
                                canvas.width = targetWidth;
                                canvas.height = targetHeight;
                                // 绘制并缩小
                                ctx.drawImage(video, 0, 0, w, h, 0, 0, targetWidth, targetHeight);
                                // 获取Base64格式的压缩图片
                                const compressedDataURL = canvas.toDataURL('image/jpeg', 0.8);
                                console.log(compressedDataURL);
                                // 移除DataURL前缀，只保留Base64数据
                                const base64Data = compressedDataURL.split(',')[1];
                                window?.__onVideoFrame?.({
                                    width: targetWidth,  // 压缩后的宽度
                                    height: targetHeight, // 压缩后的高度
                                    data: base64Data,  // Base64编码的图片数据
                                    ts: video.currentTime,
                                    format: 'image/jpeg',
                                    encoding: 'base64',  // 指定编码格式
                                    originalWidth: w,  // 原始宽度（可选）
                                    originalHeight: h,  // 原始高度（可选）
                                });
                            }
                        }
                        lastCaptureTime = now;
                    }

                    rafId = requestAnimationFrame(captureFrame);
                };

                // 保存状态
                this.videoMap.set(video, {
                    audioCtx,
                    src,
                    compressor,
                    onAudioMessage,
                    canvas,
                    ctx,
                    rafId
                });
                captureFrame();

            } catch (error) {
                console.warn('AudioWorklet not supported:', error);

                // 降级处理：暂停音频上下文
                await audioCtx.suspend();

                // 保存状态但不处理音频
                const canvas = document.createElement('canvas');
                const ctx = canvas.getContext('2d');
                let rafId = null;
                let lastCaptureTime = 0;
                const FPS = 1;
                const INTERVAL = 1000 / FPS;

                const captureFrame = () => {
                    if (!this.videoMap.has(video)) return;
                    const now = performance.now();
                    if (now - lastCaptureTime >= INTERVAL) {
                        if (video.readyState >= 2 && !video.paused) {
                            const w = video.videoWidth;
                            const h = video.videoHeight;
                            if (w && h) {
                                const scale = 0.2;
                                const targetWidth = Math.floor(w * scale);
                                const targetHeight = Math.floor(h * scale);
                                canvas.width = targetWidth;
                                canvas.height = targetHeight;
                                ctx.drawImage(video, 0, 0, w, h, 0, 0, targetWidth, targetHeight);
                                const compressedDataURL = canvas.toDataURL('image/jpeg', 0.8);
                                console.log(compressedDataURL);
                                const base64Data = compressedDataURL.split(',')[1];
                                window?.__onVideoFrame?.({
                                    width: targetWidth,
                                    height: targetHeight,
                                    data: base64Data,
                                    ts: video.currentTime,
                                    format: 'image/jpeg',
                                    encoding: 'base64',
                                    originalWidth: w,
                                    originalHeight: h,
                                });
                            }
                        }
                        lastCaptureTime = now;
                    }
                    rafId = requestAnimationFrame(captureFrame);
                };

                this.videoMap.set(video, {
                    audioCtx,
                    src: null,
                    compressor: null,
                    onAudioMessage: null,
                    canvas,
                    ctx,
                    rafId
                });
                captureFrame();
            }
        }

        stop(videoEl) {
            if (!videoEl) return;
            const state = this.videoMap.get(videoEl);
            if (!state) return;

            const { audioCtx, rafId } = state;

            // 停止视频捕获
            if (rafId) cancelAnimationFrame(rafId);

            // 暂停音频处理
            if (audioCtx && audioCtx.state !== 'closed') {
                audioCtx.suspend();
            }

            // 更新状态，移除回调函数
            state.rafId = null;
            if (state.compressor && state.onAudioMessage) {
                state.compressor.port.onmessage = null;
            }
        }

        destroy(videoEl) {
            if (!videoEl) return;
            const state = this.videoMap.get(videoEl);
            if (!state) return;

            const { audioCtx, src, compressor, rafId } = state;

            // 停止视频捕获
            if (rafId) cancelAnimationFrame(rafId);

            // 断开音频连接
            if (compressor) {
                compressor.disconnect();
                compressor.port.onmessage = null;
            }

            if (src) src.disconnect();

            // 关闭音频上下文
            if (audioCtx && audioCtx.state !== 'closed') {
                audioCtx.close();
            }

            // 清理canvas
            if (state.canvas) {
                state.canvas.width = 0;
                state.canvas.height = 0;
            }

            this.videoMap.delete(videoEl);
        }

        stopAll() {
            for (const videoEl of this.videoMap.keys()) {
                this.stop(videoEl);
            }
        }

        destroyAll() {
            for (const videoEl of Array.from(this.videoMap.keys())) {
                this.destroy(videoEl);
            }
        }
    }

    window.__MediaCaptureController = new MediaCaptureController();
})();