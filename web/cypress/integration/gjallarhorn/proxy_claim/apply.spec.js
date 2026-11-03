const { visitAndWaitForTable } = require('../../../support/login')

describe('Proxy Claim Apply Dialog', () => {
  beforeEach(() => {
    visitAndWaitForTable('/#/layout/gjallarhorn/proxyClaim')
    cy.get('.el-table__body-wrapper', { timeout: 15000 }).should('be.visible')
    cy.get('.action-buttons .el-button').contains('代理申请').click()
    cy.get('.el-dialog').should('be.visible')
  })

  it('should display apply form with correct fields', () => {
    cy.get('.el-dialog__title').should('contain', '提交代理申请')
    cy.get('input[placeholder*="OA标识"]').should('be.visible')
    cy.get('input[placeholder="请输入申请人姓名"]').should('be.visible')
    cy.get('textarea[placeholder="请输入源项目描述"]').should('be.visible')
    cy.get('textarea[placeholder="请输入目标项目描述"]').should('be.visible')
  })

  it('should validate OA ID format', () => {
    cy.get('input[placeholder*="OA标识"]').type('invalid oa id!')
    cy.get('input[placeholder="请输入申请人姓名"]').type('test user')
    cy.get('textarea[placeholder="请输入目标项目描述"]').type('test project')
    cy.get('.el-button--primary').contains('确 定').click()
    cy.get('.el-form-item__error').should('contain', 'OA ID')
  })

  it('should accept valid OA ID with hyphens', () => {
    cy.get('input[placeholder*="OA标识"]').type('OA-TEST-E2E-001')
    cy.get('input[placeholder="请输入申请人姓名"]').type('test user')
    cy.get('textarea[placeholder="请输入目标项目描述"]').type('test project')
    cy.get('.el-button--primary').contains('确 定').click()
    cy.get('.el-form-item__error').should('not.exist')
  })

  it('should close dialog when clicking cancel', () => {
    cy.get('.el-dialog .dialog-footer .el-button').contains('取 消').click({ force: true })
    cy.wait(300)
    cy.get('.el-dialog__wrapper').should('not.be.visible')
  })

  it('should reset form when dialog closes and reopens', () => {
    cy.get('.el-dialog .dialog-footer .el-button').contains('取 消').click({ force: true })
    cy.wait(300)
    cy.get('.el-dialog__wrapper').should('not.be.visible')
    
    cy.get('.action-buttons .el-button').contains('代理申请').click()
    cy.get('.el-dialog').should('be.visible')
    cy.get('input[placeholder*="OA标识"]').should('have.value', '')
    cy.get('input[placeholder="请输入申请人姓名"]').should('have.value', '')
  })
})