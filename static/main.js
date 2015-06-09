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
    optionsEl.innerHTML += '<li><img src="' + gif + '" /></li>'
  })
}

Gtcha.prototype.onCheckboxClick = function (e) {
  e.preventDefault()
  this._el.querySelector('label').innerHTML = 'Prove it by selecting gifs w/ ' + this._data.tag
  this._el.className += ' active'
}

Gtcha.prototype.onGifSelect = function (e) {}
