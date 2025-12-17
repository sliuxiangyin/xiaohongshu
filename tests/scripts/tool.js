// 获取当前滚动距离
function smoothScrollTo(element, distance, duration = 1000) {
    const target = typeof element === 'string' ? document.querySelector(element) : element;
    if (!target) return;

    const start = window.pageYOffset || document.documentElement.scrollTop;
    const targetPosition = start + distance;
    let startTime = null;

    function animation(currentTime) {
        if (startTime === null) startTime = currentTime;
        const timeElapsed = currentTime - startTime;
        const progress = Math.min(timeElapsed / duration, 1);

        // 使用缓动函数
        const easeProgress = easeInOutCubic(progress);
        const currentPosition = start + distance * easeProgress;

        window.scrollTo(0, currentPosition);

        if (timeElapsed < duration) {
            requestAnimationFrame(animation);
        }
    }

    function easeInOutCubic(t) {
        return t < 0.5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2;
    }

    requestAnimationFrame(animation);
}

function getScrollDistance(){
    return {
        scrollY: window.pageYOffset || document.documentElement.scrollTop || document.body.scrollTop,
        scrollX: window.pageXOffset || document.documentElement.scrollLeft || document.body.scrollLeft
    };
}
// 获取元素可见高度（视口内可见部分）+ 滚动距离
function getVisibleHeight(element) {
    const rect = element.getBoundingClientRect();
    const windowHeight = window.innerHeight || document.documentElement.clientHeight;
    const scrollY = getScrollDistance().scrollY;

    const visibleTop = Math.max(rect.top, 0);
    const visibleBottom = Math.min(rect.bottom, windowHeight);
    const visibleHeight = Math.max(visibleBottom - visibleTop, 0);

    // 元素在文档中的实际位置（加上滚动距离）
    const absoluteTop = rect.top + scrollY;
    const absoluteBottom = rect.bottom + scrollY;

    return {
        visibleHeight: visibleHeight,
        visibleTop: visibleTop,
        visibleBottom: visibleBottom,
        absoluteTop: absoluteTop,
        absoluteBottom: absoluteBottom,
        scrollY: scrollY,
        isFullyVisible: rect.top >= 0 && rect.bottom <= windowHeight,
        isPartiallyVisible: rect.top < windowHeight && rect.bottom >= 0
    };
}
// 获取元素实际高度（包括padding、border）+ 滚动距离
function getActualHeight(element) {
    const scrollY = getScrollDistance().scrollY;
    const rect = element.getBoundingClientRect();

    return {
        actualHeight: element.offsetHeight,
        top: rect.top + scrollY,
        bottom: rect.bottom + scrollY,
        scrollY: scrollY
    };
}
// 获取元素完整信息
function getElementInfo(element) {
    const scrollInfo = getScrollDistance();
    const rect = element.getBoundingClientRect();
    const windowHeight = window.innerHeight;

    const visibleTop = Math.max(rect.top, 0);
    const visibleBottom = Math.min(rect.bottom, windowHeight);
    const visibleHeight = Math.max(visibleBottom - visibleTop, 0);

    return {
        scroll: scrollInfo,
        element: {
            // 高度信息
            visibleHeight: visibleHeight,
            actualHeight: element.offsetHeight,
            contentHeight: element.scrollHeight,

            // 相对视口位置
            viewportTop: rect.top,
            viewportBottom: rect.bottom,

            // 绝对位置（文档中的位置）
            absoluteTop: rect.top + scrollInfo.scrollY,
            absoluteBottom: rect.bottom + scrollInfo.scrollY,

            // 可见性状态
            isFullyVisible: rect.top >= 0 && rect.bottom <= windowHeight,
            isPartiallyVisible: rect.top < windowHeight && rect.bottom >= 0,
            visibleRatio: visibleHeight / element.offsetHeight,

            // 边界信息
            topFromViewportTop: rect.top,
            bottomFromViewportBottom: windowHeight - rect.bottom
        }
    };
}

