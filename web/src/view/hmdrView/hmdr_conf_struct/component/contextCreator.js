export class NewContextGenerator {
  constructor(data) {
    this.data = {}
  }
  generateNewContextOptions() {
    return this.data
  }
}

export class _ContextCreator {
  constructor() {
    this.dialogVisible = false
    this.defaultTagNodeProps = {
      label: 'label',
      isLeaf: 'isLeaf'
    }
    this.formData = this.defaultFormData()
    this.tagNodeData = this.defaultTagNodeData()
  }
  get camelCaseCtxName() { return 'Context' }
  className() {
    return this.camelCaseCtxName.charAt(0).toUpperCase() + this.camelCaseCtxName.slice(1) + 'Creator'
  }
  elementName() {
    return this.className().replace(/([a-z])([A-Z])/g, '$1-$2').toLowerCase()
  }
  contextName() {
    return this.camelCaseCtxName.replace(/([a-z])([A-Z])/g, '$1-$2').toLowerCase()
  }
  defaultTagNodeData() {
    return [
      {
        id: 'new-' + this.contextName(),
        label: this.contextName(),
        ctxType: this.contextName(),
        className: this.className(),
        dragType: 'create',
        isLeaf: false
      }
    ]
  }
  resetTagNodeData() { this.tagNodeData = this.defaultTagNodeData() }
  tagTreeRefName() { return this.camelCaseCtxName + 'TagTree' }
  dialogTitle() { return '新建' + this.camelCaseCtxName }
  formRefName() { return this.className() + 'From' }
  defaultFormData() {
    return {}
  }
  resetFormData() { this.formData = this.defaultFormData() }
  formRules() { return {} }
  initFormDataWithTargetNode(targetNode, position) {}
  openDialog() { this.dialogVisible = true }
  closeDialog() { this.dialogVisible = false }
  generateSubmitOptions() {
    return new NewContextGenerator(this.formData).generateNewContextOptions()
  }
}
