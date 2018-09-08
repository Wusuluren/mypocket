//index.js
//获取应用实例
const app = getApp()

var serverInfo = require('./serverInfo.js');
var gThis = undefined

Page({
  data: {
    serverUrl: "",
    serverUser: "",
    serverPasswd: "",
    hasServerInfo: false,
    bookmark: '',
    item: 40,
    sendSuccess: true
  },
  //事件处理函数
  onLoad: function () {
    console.log(serverInfo)
    gThis = this
    var info = serverInfo.serverInfo
    if (serverInfo !== undefined) {
      this.serverUrl = info.url
      this.serverUser = info.user
      this.serverPasswd = info.passwd
      this.hasServerInfo = true
      this.setData({
        hasServerInfo: true
      })
    }
  },
  bindBookmarkInput: function(e) {
    this.bookmark = e.detail.value
  },
  bindUrlInput: function (e) {
    this.serverUrl = e.detail.value
  },
  bindUserInput: function (e) {
    this.serverUser = e.detail.value
  },
  bindPasswdInput: function (e) {
    this.serverPasswd = e.detail.value
  },
  doPost: function(action) {
    console.log(action, this.bookmark)
    wx.request({
      url: this.serverUrl + '/' + action,
      method: "POST",
      data: {
        "user": this.serverUser,
        "passwd": this.serverPasswd,
        "title": this.bookmark,
        "url": this.bookmark
      },
      header: {
        'content-type': 'application/json'
      },
      success: function (res) {
        this.sendSuccess = true
        gThis.setSendStatus(true)
      },
      fail: function(res) {
        this.sendSuccess = false
        gThis.setSendStatus(false)
      }
    })
  },
  setSendStatus: function(status) {
    this.setData({
      sendSuccess: status
    })
  },
  onAdd: function (e) {
    this.doPost('add')
  },
  onRem: function (e) {
    this.doPost('del')
  },
  setServerInfo: function(e) {
    this.setData({
      hasServerInfo: true
    })
  }
})
