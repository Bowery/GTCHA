/* PLEASE DO NOT COPY AND PASTE THIS CODE. */
var gtcha = {}
window.onload = function () {
  var gtchaEl = document.querySelector('.gtcha')
  var siteKey = gtchaEl.getAttribute('data-key')
  var iframe = document.createElement('iframe')
  iframe.style.border = 'none'
  iframe.style.transition = '200ms height linear'
  iframe.src = '/static/gtcha.html?site_key=' + siteKey
  gtchaEl.appendChild(iframe)

  var success = false

  /**
   * getResponse gets the status of the GTCHA.
   * Returns true if the user passes, false
   * if the user fails.
   * @return {Boolean}
   */
  gtcha.getResponse = function () {
    return success
  }

  /**
   * In the event the GTCHA needs to be reset,
   * the user may do so.
   */
  gtcha.reset = function () {
    success = false
    iframe.contentWindow.postMessage('reset', '*')
  }

  window.addEventListener('message', function (e) {
    switch (e.data[0]) {
      case 'setHeight':
        iframe.style.height = e.data[1] + 'px'
      case 'setResponse':
        success = e.data[1]
    }
  })
}
