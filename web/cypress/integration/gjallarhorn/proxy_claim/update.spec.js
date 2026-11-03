const { visitAndWaitForTable } = require('../../../support/login')

describe('Proxy Claim Update Dialog', () => {
  beforeEach(() => {
    visitAndWaitForTable('/#/layout/gjallarhorn/proxyClaim')
  })

  it('should open update dialog when clicking update button', () => {
    cy.get('.el-table__row').contains('更新').parents('.el-table__row').first().find('.el-button').contains('更新').click()
    cy.get('.el-dialog__title').should('contain', '更新代理申请')
    cy.get('.el-dialog').should('be.visible')
    cy.get('.el-dialog .dialog-footer .el-button').contains('取 消').click({ force: true })
    cy.wait(300)
  })

  it('should display existing data in update form', () => {
    cy.get('.el-table__row').contains('更新').parents('.el-table__row').first().find('.el-button').contains('更新').click()
    cy.get('.el-dialog').should('be.visible')
    cy.get('input[disabled]').invoke('val').should('not.be.empty')
    cy.get('.el-dialog .dialog-footer .el-button').contains('取 消').click({ force: true })
    cy.wait(300)
  })

  it('should show OA ID as disabled in update mode', () => {
    cy.get('.el-table__row').contains('更新').parents('.el-table__row').first().find('.el-button').contains('更新').click()
    cy.get('input[placeholder*="OA标识"]').should('be.disabled')
    cy.get('.el-dialog .dialog-footer .el-button').contains('取 消').click({ force: true })
    cy.wait(300)
  })

  it('should hide applicant field in update mode', () => {
    cy.get('.el-table__row').contains('更新').parents('.el-table__row').first().find('.el-button').contains('更新').click()
    cy.get('input[placeholder="请输入申请人姓名"]').should('not.exist')
    cy.get('.el-dialog .dialog-footer .el-button').contains('取 消').click({ force: true })
    cy.wait(300)
  })

  it('should close and reopen update dialog correctly', () => {
    cy.get('.el-table__row').contains('更新').parents('.el-table__row').first().find('.el-button').contains('更新').click()
    cy.get('.el-dialog').should('be.visible')
    cy.get('.el-dialog .dialog-footer .el-button').contains('取 消').click({ force: true })
    cy.wait(300)
    cy.get('.el-dialog__wrapper').should('not.be.visible')

    cy.get('.el-table__row').contains('更新').parents('.el-table__row').first().find('.el-button').contains('更新').click()
    cy.get('.el-dialog__title').should('contain', '更新代理申请')
    cy.get('.el-dialog .dialog-footer .el-button').contains('取 消').click({ force: true })
    cy.wait(300)
  })
})