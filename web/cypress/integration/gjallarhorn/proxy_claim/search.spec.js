const { visitAndWaitForTable } = require('../../../support/login')

describe('Proxy Claim Page', () => {
  beforeEach(() => {
    visitAndWaitForTable('/#/layout/gjallarhorn/proxyClaim')
    cy.get('.el-table__body-wrapper', { timeout: 15000 }).should('be.visible')
  })

  it('should display search form and table', () => {
    cy.get('.search-term').should('be.visible')
    cy.get('.el-table').should('be.visible')
  })

  it('should open apply dialog when clicking 代理申请 button', () => {
    cy.get('.action-buttons .el-button').contains('代理申请').click()
    cy.get('.el-dialog__title').should('contain', '提交代理申请')
    cy.get('.el-dialog').should('be.visible')
  })
})