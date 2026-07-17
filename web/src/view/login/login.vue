<template>
  <div id="userLayout" class="user-layout-wrapper">
    <div class="container">
      <div class="top">
        <div class="desc">
          <img class="logo_login" src="@/assets/logo_login.png" alt="">
        </div>
        <div class="header" />
      </div>
      <div class="main">
        <el-form
          ref="loginForm"
          :model="loginForm"
          :rules="rules"
          @keyup.enter.native="submitForm"
        >
          <el-form-item prop="username">
            <el-input
              v-model="loginForm.username"
              placeholder="请输入用户名"
            >
              <i
                slot="suffix"
                class="el-input__icon el-icon-user"
              /></el-input>
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              :type="lock === 'lock' ? 'password' : 'text'"
              placeholder="请输入密码"
            >
              <i
                slot="suffix"
                :class="'el-input__icon el-icon-' + lock"
                @click="changeLock"
              />
            </el-input>
          </el-form-item>
          <el-form-item style="position:relative">
            <el-input
              v-model="loginForm.captcha"
              name="logVerify"
              placeholder="请输入验证码"
              style="width:60%"
            />
            <div class="vPic">
              <img
                v-if="picPath"
                :src="picPath"
                width="100%"
                height="100%"
                alt="请输入验证码"
                @click="loginVerify()"
              >
            </div>
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              style="width:100%"
              @click="submitForm"
            >登 录</el-button>
          </el-form-item>
        </el-form>
      </div>

      <div class="footer">
        <div class="copyright">
          Copyright &copy; {{ curYear }} 💖测试环境组
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapActions } from 'vuex'
import { captcha, getPublicKey } from '@/api/user'
import JSEncrypt from 'jsencrypt/bin/jsencrypt.min'

export default {
  name: 'Login',
  data() {
    const checkUsername = (rule, value, callback) => {
      if (value.length < 5 || value.length > 12) {
        return callback(new Error('请输入正确的用户名'))
      } else {
        callback()
      }
    }
    const checkPassword = (rule, value, callback) => {
      if (value.length < 6 || value.length > 12) {
        return callback(new Error('请输入正确的密码'))
      } else {
        callback()
      }
    }
    return {
      curYear: 0,
      lock: 'lock',
      loginForm: {
        username: 'admin',
        password: '123456',
        captcha: '',
        captchaId: ''
      },
      rules: {
        username: [{ validator: checkUsername, trigger: 'blur' }],
        password: [{ validator: checkPassword, trigger: 'blur' }]
      },
      logVerify: '',
      picPath: '',
      publicKey: '',
      challenge: ''
    }
  },
  created() {
    this.curYear = new Date().getFullYear()
    this.initLogin()
  },
  methods: {
    ...mapActions('user', ['LoginIn']),
    async initLogin() {
      // Get captcha and public key with challenge
      await this.loginVerify()
    },
    async fetchPublicKey() {
      try {
        const res = await getPublicKey({ captchaId: this.loginForm.captchaId })
        if (res.code === 0) {
          this.publicKey = res.data.publicKey
          this.challenge = res.data.challenge
        }
      } catch (error) {
        console.error('Failed to fetch public key:', error)
      }
    },
    encryptLoginData(loginData) {
      if (!this.publicKey || !this.challenge) {
        return null
      }
      const encryptor = new JSEncrypt()
      encryptor.setPublicKey(this.publicKey)

      // Add challenge for replay attack prevention
      const dataWithChallenge = {
        ...loginData,
        challenge: this.challenge
      }

      // Encrypt the entire JSON string
      const jsonString = JSON.stringify(dataWithChallenge)
      return encryptor.encrypt(jsonString)
    },
    async login() {
      const encryptedData = this.encryptLoginData(this.loginForm)
      if (!encryptedData) {
        this.$message({
          type: 'error',
          message: '加密初始化失败，请刷新页面',
          showClose: true
        })
        return
      }
      await this.LoginIn({ encrypted_data: encryptedData })
    },
    async submitForm() {
      this.$refs.loginForm.validate(async(v) => {
        if (v) {
          this.login()
          this.loginVerify()
        } else {
          this.$message({
            type: 'error',
            message: '请正确填写登录信息',
            showClose: true
          })
          this.loginVerify()
          return false
        }
      })
    },
    changeLock() {
      this.lock === 'lock' ? (this.lock = 'unlock') : (this.lock = 'lock')
    },
    loginVerify() {
      return captcha({}).then((ele) => {
        this.picPath = ele.data.picPath
        this.loginForm.captchaId = ele.data.captchaId
        // After getting captchaId, fetch public key and challenge
        return this.fetchPublicKey()
      })
    }
  }
}
</script>

<style scoped lang="scss">
@import '@/style/login.scss';

</style>
