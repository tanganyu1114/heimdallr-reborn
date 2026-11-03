import './commands'

// 抑制 Vue Router 重定向错误
Cypress.on('uncaught:exception', (err, runnable) => {
  if (err.message.includes('Redirected') || err.message.includes('NavigationDuplicated')) {
    return false
  }
  return true
})