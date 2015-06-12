/* PLEASE DO NOT COPY AND PASTE THIS CODE. */
window.onload = function () {
  var gtchaEl = document.querySelector('.gtcha')
  var siteKey = gtchaEl.getAttribute('data-key')
  var iframe = document.createElement('iframe')
  iframe.style.border = 'none'
  iframe.style.transition = '200ms height linear'
  iframe.src = '/static/gtcha.html?site_key=' + siteKey
  gtchaEl.appendChild(iframe)

  window.addEventListener('message', function (e) {
    if (e.data[0] == 'setHeight') {
      iframe.style.height = e.data[1] + 'px'
    }
  })
}
