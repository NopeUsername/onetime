/* Setup Ace */
const editor = ace.edit("ace")
editor.setOptions({
    highlightActiveLine: true,
    highlightSelectedWord: true,
    animatedScroll: true,
    showPrintMargin: false,
    theme: "ace/theme/monokai",
    enableBasicAutocompletion: true,
    enableLiveAutocompletion: true,
    behavioursEnabled: true,

    fontSize: 13,
    tabSize: 4,

    mode: "ace/mode/lua"
})

/* Utilities */
function makeNotification(type, text) {
    const notif = $(`<div class="tile notif-container is-parent is-vertical is-4">
        <div class="notification ${type}">
            ${text}
            <button class="delete"></button>
        </div>
    </div>`)
    $(document.body).append(notif)
    
    notif.find("button").on("click", () => {
        notif.remove()
    })

    setTimeout(() => {
        notif.remove()
    }, 10000)
}

function copyToClipboard(text) {
    const el = document.createElement("textarea")
    document.body.appendChild(el)
    el.value = text
    el.select()

    document.execCommand("copy")
    document.body.removeChild(el)
}

/* Uploading script */
const path = window.location.href

$("#upload").on("click", () => {
    const value = editor.getValue().replace(/\s/g, "")
    
    if (value.length == 0) {
        return makeNotification("is-danger", "You have not entered any text.")
    }
    $("#upload").prop("disabled", true)

    const xhr = new XMLHttpRequest()

    xhr.addEventListener("readystatechange", () => {
        if (xhr.readyState == xhr.DONE) {
            $("#upload").prop("disabled", false)

            if (xhr.status == 200) {
                makeNotification("is-success", "Copied URL to clipboard.")
                copyToClipboard(location + "p/" + xhr.responseText)
            } else {
                makeNotification("is-danger", `Error (${xhr.status}): ${xhr.responseText}`)
            }
        }
    })
    xhr.open("POST", "/upload")
    xhr.send(value)
})