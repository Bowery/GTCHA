window.onload = function() {
  var formEl = document.querySelector('form')
  var name = ''
  var domains = ''

  formEl.onsubmit = function(e) {
    e.preventDefault()
    name = document.querySelector('form #name').value
    domains = document.querySelector('form #domains').value.split(',')
    domains = 'domain=' + domains.map(function(el) {
      return el.trim()
    }).join('&domain=')
    var params = 'name=' + name + '&' + domains

    var req = new XMLHttpRequest()
    req.open('POST', window.location.origin + '/register')
    req.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
    req.onload = function() {
      if (req.status !== 200)
        return

      var data = JSON.parse(req.responseText)
      showEmbedCode(data.api_key)
    }

    req.send(params)
  }
}

/**
 * showEmbedCode displays the embed code with
 * the proper api key.
 * @param {String} key
 */
function showEmbedCode (key) {
  var wrapperEl = document.querySelector('.embed-code')
  var embedScript = wrapperEl.querySelector('.embed-script')
  var embedHTML = wrapperEl.querySelector('.embed-html')
  embedScript.value = '<script src="' + window.location.origin + '/static/api.js" async defer></script>'
  embedScript.onclick = function (e) {
    e.target.select()
  }
  embedHTML.value = '<div class="gtcha data-key="' + key + '"></div>'
  embedHTML.onclick = function (e) {
    e.target.select()
  }

  wrapperEl.className = 'embed-code active'
}
