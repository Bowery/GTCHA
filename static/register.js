window.onload = function() {
  var formEl = document.querySelector("form")
  var name = ""
  var domains = ""

  formEl.onsubmit = function(e) {
    e.preventDefault()
    name = document.querySelector("form #name").value
    domains = document.querySelector("form #domains").value.split(",")
    domains = "domain=" + domains.map(function(el) {
      return el.trim()
    }).join("&domain=")
    var params = "name=" + name + "&" + domains

    var req = new XMLHttpRequest()
    req.open("POST", "http://localhost:8080/register")
    req.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    req.onload = function() {
      if (req.status !== 200)
        return

      var data = JSON.parse(req.responseText)
      console.log(data)
    }

    req.send(params)
  }
}
