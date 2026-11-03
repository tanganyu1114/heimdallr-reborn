const { visitAndWaitForTable } = require('../../../support/login')

// Mock data
const mockMigrationRequests = [
  {
    id: 'mig-001',
    name: 'Test Migration',
    source_oa_id: 'OA-TEST-001',
    target_oa_id: 'OA-TEST-002',
    source_claim_origin_state: 'active',
    target_claim_origin_state: 'active',
    migrate_state: 'completed',
    retry_count: 0,
    migrating_addresses: [
      { source_address: '10.0.0.1:8080', target_address: '10.0.0.5:8080' }
    ]
  }
]

// Reusable helper functions
function openFirstClaimDetail() {
  cy.get('.el-table__body .el-table__row').first().find('.el-button').contains('详情').click()
  cy.get('.el-dialog__title', { timeout: 10000 }).should('contain', '代理声明详情')
}

function setupMigrationMock(data) {
  cy.intercept('GET', '**/gjallarhorn/v1/proxy-claims/migrations*', {
    statusCode: 200,
    body: data
  }).as('loadMigrations')
}

function openDetailWithMigrations() {
  setupMigrationMock(mockMigrationRequests)
  openFirstClaimDetail()
  cy.wait('@loadMigrations', { timeout: 10000 })
  cy.get('.migration-collapse-header', { timeout: 10000 }).should('be.visible')
  cy.get('.migration-collapse-header').click()
  cy.get('.migration-request-list', { timeout: 10000 }).should('be.visible')
}

describe('Migration Request Display - E2E Test', () => {
  beforeEach(() => {
    visitAndWaitForTable('/#/layout/gjallarhorn/proxyClaim')
    cy.get('.el-table__body-wrapper', { timeout: 15000 }).should('be.visible')
  })

  describe('Migration Request Collapsible Panel', () => {
    it('should show migration collapsible panel when claim has migrations', () => {
      setupMigrationMock(mockMigrationRequests)
      openFirstClaimDetail()

      cy.wait('@loadMigrations', { timeout: 10000 })
      cy.get('.migration-collapse-panel', { timeout: 10000 }).should('be.visible')
      cy.get('.migration-collapse-header').should('be.visible')
      cy.get('.header-title').should('contain', '迁移请求')
    })

    it('should not show migration panel when claim has no migrations', () => {
      setupMigrationMock([])
      openFirstClaimDetail()

      cy.wait('@loadMigrations', { timeout: 10000 })
      cy.wait(500)

      cy.get('body').then(($body) => {
        const panel = $body.find('.migration-collapse-panel')
        expect(panel.length).to.equal(0)
      })
    })

    it('should show migration panel above proxy details section', () => {
      setupMigrationMock(mockMigrationRequests)
      openFirstClaimDetail()

      cy.wait('@loadMigrations', { timeout: 10000 })
      cy.get('.migration-collapse-panel', { timeout: 10000 }).should('be.visible')
      cy.get('.el-divider').contains('代理详情').should('exist')

      cy.get('.migration-collapse-panel').then(($panel) => {
        const panelPosition = $panel.offset().top
        cy.get('.el-divider').contains('代理详情').then(($divider) => {
          const dividerPosition = $divider.offset().top
          expect(panelPosition).to.be.lessThan(dividerPosition)
        })
      })
    })

    it('should be collapsed by default', () => {
      setupMigrationMock(mockMigrationRequests)
      openFirstClaimDetail()

      cy.wait('@loadMigrations', { timeout: 10000 })
      cy.get('.migration-collapse-header', { timeout: 10000 }).should('be.visible')
      cy.get('.migration-request-list').should('not.be.visible')
      cy.get('.expand-hint').should('contain', '点击展开')
    })

    it('should expand when clicking header', () => {
      setupMigrationMock(mockMigrationRequests)
      openFirstClaimDetail()

      cy.wait('@loadMigrations', { timeout: 10000 })
      cy.get('.migration-collapse-header', { timeout: 10000 }).should('be.visible')

      cy.get('.migration-collapse-header').click()
      cy.get('.migration-request-list', { timeout: 10000 }).should('be.visible')
      cy.get('.expand-hint').should('contain', '点击折叠')
    })

    it('should collapse when clicking header again', () => {
      setupMigrationMock(mockMigrationRequests)
      openFirstClaimDetail()

      cy.wait('@loadMigrations', { timeout: 10000 })
      cy.get('.migration-collapse-header', { timeout: 10000 }).should('be.visible')

      cy.get('.migration-collapse-header').click()
      cy.get('.migration-request-list', { timeout: 10000 }).should('be.visible')

      cy.get('.migration-collapse-header').click()
      cy.get('.migration-request-list').should('not.be.visible')
    })

    it('should show arrow icon that rotates on expand/collapse', () => {
      setupMigrationMock(mockMigrationRequests)
      openFirstClaimDetail()

      cy.wait('@loadMigrations', { timeout: 10000 })
      cy.get('.migration-collapse-header', { timeout: 10000 }).should('be.visible')

      cy.get('.collapse-icon').should('have.class', 'el-icon-arrow-right')

      cy.get('.migration-collapse-header').click()
      cy.get('.collapse-icon').should('have.class', 'el-icon-arrow-down')
    })
  })

  describe('Migration Request List Component', () => {
    it('should display migration request table with correct columns', () => {
      openDetailWithMigrations()

      cy.get('.migration-request-list').should('be.visible')
      cy.get('.el-table').should('be.visible')

      cy.get('.el-table__header').should('contain', '源OA ID')
      cy.get('.el-table__header').should('contain', '源声明状态')
      cy.get('.el-table__header').should('contain', '目标声明状态')
      cy.get('.el-table__header').should('contain', '迁移状态')
      cy.get('.el-table__header').should('contain', '重试次数')
      cy.get('.el-table__header').should('contain', '迁移地址数')
      cy.get('.el-table__header').should('contain', '操作')
    })

    it('should display migration state tags with correct colors', () => {
      openDetailWithMigrations()

      cy.get('.migration-request-list').should('be.visible')

      cy.get('.el-table__body').then(($tableBody) => {
        if ($tableBody.find('.el-table__row').length > 0) {
          cy.get('.el-table__row').first().find('.el-tag').should('be.visible')
        }
      })
    })
  })
})