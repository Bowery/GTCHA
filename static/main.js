/* PLEASE DO NOT COPY AND PASTE THIS CODE. */
window.onload = function () {
  gtcha = new Gtcha()
  gtcha.init()
}

function Gtcha () {}

Gtcha.prototype._el = null

Gtcha.prototype._data = null

Gtcha.prototype.init = function () {
  this._el = document.querySelector('.gtcha')
  this._el.className = 'gtcha'
  this._el.querySelector('input[type="checkbox"]').onclick = this.onCheckboxClick.bind(this)
  this.fetch()
}

Gtcha.prototype.fetch = function () {
  var req = new XMLHttpRequest()
  var self = this
  req.onreadystatechange = function () {
    if (req.readyState === 4) {
      self.onResponse(req)
    }
  }
  req.open('GET', '/dummy', true)
  req.send('')
}

Gtcha.prototype.onResponse = function (res) {
  if (res.status != 200) {
    // handle error
  }

  this._data = JSON.parse(res.response)
  var optionsEl = this._el.querySelector('.options')
  var self = this
  this._data.images.forEach(function (gif) {
    optionsEl.innerHTML += '\
      <div class="option">\
        <div class="mask"></div>\
        <img src="' + gif + '">\
      </div>'
  })
}

Gtcha.prototype.onCheckboxClick = function (e) {
  e.preventDefault()

  // Update desc display.
  this._el.querySelector('label').innerHTML = this._data.tag

  // Expand height to accomodate 
  this._el.className += ' active'
  this._el.style.transitionDelay = '0s'
}

Gtcha.prototype.onGifSelect = function (e) {}
