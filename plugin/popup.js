function hostClick() {
    var host = document.getElementById("host").value
    chrome.storage.local.set({'host':host})
}
function userClick() {
    var user = document.getElementById("user").value
    chrome.storage.local.set({'user':user})
}
function passwdClick() {
    var passwd = document.getElementById("passwd").value
    chrome.storage.local.set({'passwd':passwd})
}

window.onload = function () {
    chrome.storage.local.get(null, function (config) {
        if (config.host != undefined) {
            document.getElementById("host").value = config.host
        }
        if (config.user != undefined) {
            document.getElementById("user").value = config.user
        }
        if (config.passwd != undefined) {
            document.getElementById("passwd").value = config.passwd
        }
    })

    document.getElementById("host-btn").addEventListener('click', hostClick, true)
    document.getElementById("user-btn").addEventListener('click', userClick, true)
    document.getElementById("passwd-btn").addEventListener('click', passwdClick, true)
}