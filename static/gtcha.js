window.onload = function () {
  gtcha = new Gtcha()
  gtcha.init()
}

/**
 * @constructor
 */
function Gtcha () {}

/**
 * _el represents the GTCHA DOM element.
 */
Gtcha.prototype._el = null

/**
 * _data represents the GTCHA data.
 */
Gtcha.prototype._data = null

/**
 * _key represents the site key.
 */
Gtcha.prototype._key = ''

/**
 * init initiates the animation, sets click event
 * for the checkbox, and fetches data.
 */
Gtcha.prototype.init = function () {
  this._el = document.querySelector('.gtcha')
  this._el.className = 'gtcha'
  this._key = window.location.href.split('?site_key=')[1]

  this.updateHeight(84)
  this.fetch()
  window.addEventListener('message', this.onResetRequest.bind(this))

  this._el.querySelector('input[type="checkbox"]').onclick = this.onCheckboxClick.bind(this)
  this._el.querySelector('form').onsubmit = this.onSubmitClick.bind(this)
}

/**
 * fetch executs an AJAX request to the GTCHA service,
 * retrieving the verification prompt and options.
 */
Gtcha.prototype.fetch = function () {
  var req = new XMLHttpRequest()
  var self = this
  req.onreadystatechange = function () {
    if (req.readyState == 4) {
      self.onResponse(req)
    }
  }
  req.open('GET', '/dummy_get', true)
  req.send('')
}

/**
 * onResponse handles the fetch response. If the response
 * is 200 OK, render the options and set click handlers.
 * @param {Object} res
 */
Gtcha.prototype.onResponse = function (res) {
  if (res.status != 200) {
    this.handleErr(JSON.parse(res.response).error)
    return
  }
  this._data = JSON.parse(res.response)
  
  var optionsEl = this._el.querySelector('.options')
  for (var i = 0; i < this._data.images.length; i++) {
    var gif = this._data.images[i]
    optionsEl.innerHTML += '\
      <div class="option" data-id="' + gif.id + '" style="background-image: url(' + gif.uri + ')"></div>'
  }

  var options = this._el.querySelectorAll('.option')
  for (var i = 0; i < options.length; i++) {
    options[i].onclick = this.onGifSelect.bind(this)
  }
}

/**
 * onCheckboxClick handles the checkbox click event.
 * It updates the prompt and displays options.
 * @param {MouseEvent}
 */
Gtcha.prototype.onCheckboxClick = function (e) {
  e.preventDefault()
  this._el.querySelector('label').innerHTML = this._data.tag
  this._el.className = 'gtcha active'
  this._el.style.transitionDelay = '0s'
  this.updateHeight(376)
}

/**
 * onGifSelect toggles the 
 * @param {MouseEvent}
 */
Gtcha.prototype.onGifSelect = function (e) {
  var el = e.target
  el.className == 'option selected'
  ? el.className = 'option'
  : el.className = 'option selected'
}

/**
 * onSubmitClick posts the users response to the GTCHA service.
 * @param {MouseEvent}
 */
Gtcha.prototype.onSubmitClick = function (e) {
  e.preventDefault()

  var optionEls = document.querySelectorAll('.option.selected')
  var options = []
  for (var i = 0; i < optionEls.length; i++) {
    options[i] = optionEls[i].getAttribute('data-id')
  }

  e.target.querySelector('input[type="submit"]').value = 'submitting...'

  var payload = {
    id: this._data.id,
    tag: this._data.tag,
    in: options
  }
  var req = new XMLHttpRequest()
  var self = this
  req.onreadystatechange = function () {
    if (req.readyState == 4) {
      self.onSubmitResponse(req)
    }
  }
  req.open('PUT', '/dummy_put?api_key=' + this._key, true)
  req.setRequestHeader('Content-Type', 'application/json')
  req.send(JSON.stringify(payload))
}

/**
 * onSubmitResponse handles the submit response.
 * @param {Object}
 */
Gtcha.prototype.onSubmitResponse = function (res) {
  if (res.status != 200) {
    this.handleErr(JSON.parse(res.response).error)
    return
  }


  this._el.className = 'gtcha'
  this._el.querySelector('label').innerHTML = 'All good.'
  this.updateHeight(84)
  this.updateResponse(true)
}

/**
 * updateHeight updates the height of the parent iframe.
 * @param {Number} height
 */
Gtcha.prototype.updateHeight = function (height) {
  window.parent.postMessage(['setHeight', height], '*')
}

/**
 * updateResponse updates the response (true|false).
 * @param {Boolean} response
 */
Gtcha.prototype.updateResponse = function (response) {
  window.parent.postMessage(['setResponse', response], '*')
}

/**
 * onResetRequest handles reset request from parent window.
 */
Gtcha.prototype.onResetRequest = function () {
  this.reset()
}

/**
 * reset resets the GTCHA to it's initial state.
 */
Gtcha.prototype.reset = function () {
  this._data = null
  this._el.className = 'gtcha'
  this._el.querySelector('label').innerHTML = 'I\'m Human.'
  this.updateHeight(84)
}

/**
 * handleErr handles the error body from a non-200 response.
 * @param {String}
 */
Gtcha.prototype.handleErr = function (err) {
  console.log(err)
}
