package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func ScriptScrollSeperatorIntoView() Node {
	return Script(Raw(`
try {
  document.getElementById('separator').scrollIntoView({
            behavior: 'auto',
            block: 'center',
            inline: 'center'
        });
} catch (e) {
    console.error(e)
}
`))
}

func ScriptReloadPageEveryMinute() Node {
	return Script(Raw(`
console.log('starting reloader')
setInterval(() => {
    console.log('reloading page')
    window.location.reload()
    }, 60*1000)
`))
}
