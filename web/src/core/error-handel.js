import { createSysError } from '@/api/system/sysError'
import { removeInitialLoading } from '@/utils/initialLoading'

function sendErrorTip(errorInfo) {
  setTimeout(() => {
    const errorData = {
      form: errorInfo.type,
      info: `${errorInfo.message}\nStack: ${errorInfo.stack}${errorInfo.component ? `\nComponent: ${errorInfo.component.name || 'Unknown'}` : ''}${errorInfo.vueInfo ? `\nVue Info: ${errorInfo.vueInfo}` : ''}${errorInfo.source ? `\nSource: ${errorInfo.source}:${errorInfo.lineno}:${errorInfo.colno}` : ''}`,
      level: 'error',
      solution: null
    }
    
    createSysError(errorData).catch(apiErr => {
      console.error('Failed to create error record:', apiErr)
    })
  }, 0)
}
  
  window.addEventListener('unhandledrejection', (event) => {
    removeInitialLoading()
    sendErrorTip({
      type: '前端',
      message: `错误信息: ${event.reason}`,
      stack: `调用栈: ${event.reason?.stack || '没有调用栈信息'}`,
    });
  });

window.addEventListener('error', (event) => {
  removeInitialLoading()
  sendErrorTip({
    type: '前端',
    message: `错误信息: ${event.message}`,
    stack: `调用栈: ${event.error?.stack || '没有调用栈信息'}`,
    source: event.filename,
    lineno: event.lineno,
    colno: event.colno
  })
})
