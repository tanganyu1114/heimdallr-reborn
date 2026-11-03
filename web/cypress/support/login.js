/**
 * Authentication helper for headless e2e tests
 * Uses x-token header directly for API request authentication
 */
function setupAuth() {
  const e2eToken = Cypress.env('E2E_TOKEN') || 'mock-token-for-testing'
  
  // Intercept all API requests and add x-token header
  cy.intercept('/api/**', (req) => {
    req.headers['x-token'] = e2eToken
    req.headers['x-user-id'] = '1'
  })
  
  // Set Vuex state in localStorage for route guard authentication
  const vuexState = {
    user: {
      userInfo: {
        uuid: 'b35dcc93-024a-4df0-a277-9c657fbc4945',
        nickName: '超级管理员',
        headerImg: '',
        authority: '888',
        ID: 1
      },
      token: e2eToken
    }
  }
  
  cy.visit('/#/login', {
    onBeforeLoad: (win) => {
      win.localStorage.setItem('vuex', JSON.stringify(vuexState))
    }
  })
  
  // Wait for page to initialize
  cy.wait(500)
}

/**
 * Visit page with authentication setup
 * @param {string} path - The route path to visit
 */
function visit(path) {
  setupAuth()
  cy.visit(path)
  cy.get('.el-menu', { timeout: 10000 }).should('be.visible')
  cy.wait(500)
}

/**
 * Visit page with authentication and wait for page to be ready
 * @param {string} path - The route path to visit
 */
function visitAndWaitForTable(path) {
  setupAuth()
  cy.visit(path)
  // Wait for page to be ready - check for menu or table
  cy.get('.el-menu', { timeout: 10000 }).should('be.visible')
  cy.wait(1000)
}

module.exports = {
  setupAuth,
  visit,
  visitAndWaitForTable
}