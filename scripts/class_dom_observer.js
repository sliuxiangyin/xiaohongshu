(function (global) {
    /**
     * ClassChangeObserver 类：用于监听特定元素上某个 Class 的添加和移除。
     */
    class MinimalClassDOMObserver {
        constructor({
                        className = 'note-detail-mask',
                        onAdd,
                        onRemove,
                    } = {}) {
            this.className = className
            this.onAdd = onAdd
            this.onRemove = onRemove
            this._observer = null
            this._currentEl = null
        }

        start() {
            if (this._observer) return

            // 先检查一次（防止已存在）
            const existing = document.body.querySelector(`.${this.className}`)
            if (existing) {
                this._currentEl = existing
                this.onAdd?.(existing)
            }

            this._observer = new MutationObserver((mutations) => {
                for (const mutation of mutations) {
                    // 新增节点
                    mutation.addedNodes.forEach((node) => {
                        if (!(node instanceof HTMLElement)) return

                        if (
                            node.classList?.contains(this.className) ||
                            node.querySelector?.(`.${this.className}`)
                        ) {
                            const el = node.classList.contains(this.className)
                                ? node
                                : node.querySelector(`.${this.className}`)

                            if (el && el !== this._currentEl) {
                                this._currentEl = el
                                this.onAdd?.(el)
                            }
                        }
                    })

                    // 移除节点
                    mutation.removedNodes.forEach((node) => {
                        if (!(node instanceof HTMLElement)) return

                        if (
                            node === this._currentEl ||
                            node.querySelector?.(`.${this.className}`)
                        ) {
                            const removed =
                                node === this._currentEl
                                    ? node
                                    : node.querySelector?.(`.${this.className}`)

                            if (removed) {
                                this._currentEl = null
                                this.onRemove?.(removed)
                            }
                        }
                    })
                }
            })

            this._observer.observe(document.body, {
                childList: true,
                subtree: true,
            })
        }

        destroy() {
            if (!this._observer) return
            this._observer.disconnect()
            this._observer = null
            this._currentEl = null
        }
    }

    // 使用案例
    //const observer = new NoteDetailMaskObserver({
    //   onAdd(el) {
    //     console.log('note-detail-mask 被添加:', el)
    //   },
    //   onRemove(el) {
    //     console.log('note-detail-mask 被移除:', el)
    //   },
    // })
    //
    // observer.start()
    // // 销毁监听
    // // observer.destroy()
    global.MinimalClassDOMObserver = MinimalClassDOMObserver;

})(window);