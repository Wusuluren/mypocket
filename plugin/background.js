function isValidConfig(config) {
    if(config.host == undefined) {
        return false
    }
    if(config.user == undefined) {
        return false
    }
    if(config.passwd == undefined) {
        return false
    }
    return true
}

function process(action, tab) {
    chrome.storage.local.get(null, function (config) {
        if(!isValidConfig(config)) {
            alert('invalid config')
            return
        }
        var url = config.host + '/' + action
        var xhr = new XMLHttpRequest();
        xhr.open("POST", url, true)
        xhr.send(JSON.stringify({
            "user": config.user,
            "passwd": config.passwd,
            "title":tab.title,
            "url":tab.url
        }))
        xhr.onreadystatechange = function () {
            if(xhr.readyState == XMLHttpRequest.DONE) {
                if (xhr.status==200) {
                    alert('success')
                } else {
                    alert('failed:'+xhr.responseText)
                }
            }
        }
    })
}

function onclickAdd(info, tab) {
    process('add', tab)
}

function onclickDel(info, tab) {
    process('del', tab)
}

chrome.contextMenus.create({
    "title": "Save to pocket",
    "contexts": ["page"],
    "onclick": onclickAdd
})

chrome.contextMenus.create({
    "title": "Remove from pocket",
    "contexts": ["page"],
    "onclick": onclickDel
})

