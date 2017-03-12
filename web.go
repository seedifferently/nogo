package main

// nogo.css
var nogoCSS = `*,
*:after,
*:before {
  box-sizing: inherit;
}

html {
  box-sizing: border-box;
  font-size: 62.5%;
}

body {
  color: #606c76;
  font-family: 'Roboto', 'Helvetica Neue', 'Helvetica', 'Arial', sans-serif;
  font-size: 1.6em;
  font-weight: 300;
  letter-spacing: .01em;
  line-height: 1.6;
}

input[type='search'],
input[type='text'] {
  -webkit-appearance: none;
     -moz-appearance: none;
          appearance: none;
  background-color: transparent;
  border: 0.1rem solid #d1d1d1;
  border-radius: .4rem;
  box-shadow: none;
  box-sizing: inherit;
  height: 3.8rem;
  padding: .6rem 1.0rem;
  width: 100%;
}

input[type='search']:focus,
input[type='text']:focus {
  border-color: #9b4dca;
  outline: 0;
}

label {
  display: block;
  font-size: 1.6rem;
  font-weight: 700;
  margin-bottom: .5rem;
}

.container {
  margin: 0 auto;
  max-width: 112.0rem;
  padding: 0 2.0rem;
  position: relative;
  width: 100%;
}

.row {
  display: flex;
  flex-direction: column;
  padding: 0;
  width: 100%;
}

.row .column {
  display: block;
  flex: 1 1 auto;
  margin-left: 0;
  max-width: 100%;
  width: 100%;
}

@media (min-width: 40rem) {
  .row {
    flex-direction: row;
    margin-left: -1.0rem;
    width: calc(100% + 2.0rem);
  }
  .row .column {
    margin-bottom: inherit;
    padding: 0 1.0rem;
  }
}

a {
  color: #9b4dca;
  text-decoration: none;
}

a:focus, a:hover {
  color: #606c76;
}

fieldset,
input {
  margin-bottom: 1.5rem;
}

.text-right { text-align: right; }

#header {
  height: 69px;
  background-color: #2f2f2f;
  text-align: center;
}

#header a.heading {
  display: inline-block;
  float: left;
  font-weight: bold;
  font-size: 40px;
  letter-spacing: 2px;
}

#main {
  min-height: calc(100vh - 99px);
  padding-bottom: 1rem;
}

#inputs { padding-top: 15px; }

#inputs form { margin-bottom: 5px; }

#records-header > div { margin-bottom: 15px; }

#records-header,
.row.record {
  display: inline-flex;
  flex-direction: row;
  flex-wrap: nowrap;
  white-space: nowrap;
}

#back { flex: 0 0 0%; }

.row.record:hover { background-color: #eee; }

.column.actions {
  flex: 0;
  padding-right: 0;
}

.column.key { padding-left: 0; }

.actions form {
  display: inline-block;
  margin: 0;
  padding: 0;
}

.actions form.hide { display: none; }

.icon {
  display: inline-block;
  vertical-align: text-top;
  width: 16px;
  height: 16px;
  padding: 0;
  margin: 0 10px 0 0;
  border: none;
  color: inherit;
  background: none;
  cursor: pointer;
}

.icon:focus {
  outline: 0;
}

.icon-power {
  display: block;
  position: absolute;
  top: 21px;
  right: 10px;
  width: 24px;
  height: 24px;
  background: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABcAAAAYCAYAAAARfGZ1AAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAAsTAAALEwEAmpwYAAACpklEQVRIx61VPWhTURS+g8FBB4VCcXAQFLo9iqaSvKRp/kybJml+zEvSJo0xpAhBrRUVtINohKCYIQQhUCiKBMWItoMtCiWDLop0CnZwc9COuonD530vTd59L38v4nCWe8/57j3nfOc7hOM40s/MFlx5HsHWhxg2EbqfBcft0xJH+jvhcMqHrVqkYT8zKFDwE/8NPDMjg++mUYKeG/kncLsdno0wsgjeNTbS7w7O87Cuhqhv7JmZ+u7vAQ6d14n8xh7IdhxV2PhgV/CT4Ith+RzRFzfoA4c6grtsMnArYMqf7QpuQOal2j/6WHzgoBKcR/Sp2jH9JweTmH63smD4ghfrbIzEpulzaQYcwws+bCqAUz9yMDcZ0auhNHYaVTb2dxr55j0ZMyDLpicFj9v02tmC40szqqy9C4sSuNelfBm+y1eVQ9KfinoTrq+xGMnPOYoxRC4F5MOvSazCwOkH5zmOLgVln28plCnOKCkKbK2/iC8eUQUeiHnkniDzS/Q5pvLR+afk5ko0tph8ZIVNJ/E6x1KpZadhLQZQ2YnT0fecP9NpGiedqChYYzfOkUfsz5Nvc+pB0GZ0AJnefZzFOqy8QO6FGBplKI00ilKbuHllnHoCFfCcmSSYekrpOCfnBgYfg0tR3r2+EKMFhVoHGg1SkoCazkJ+WaSzeKl/GFYNQeDmotaFYBpXclxiyoTZ3dIWpw2lmlpbwk+yHZnD/Nhtw3KbeEXKt5vyS1oa4VHqi7R15in1YnU3TOkhtnkXJxEs+JWlaNS6KXZqPT8F7pZKI9SK907ofr+bopPrmDD32EQYueZt/1E/Q4r+2OEY1bBDoQs7kV0J9getxymfw2vz3QavJ8UEF6wPfLjz6izKbwRU30dQ+RRB6fss3ZvCtr53wznyF88iqrkwth+dAAAAAElFTkSuQmCC) left bottom no-repeat;
}

.icon-power[data-disabled] {
  background: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABcAAAAYCAYAAAARfGZ1AAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAAsTAAALEwEAmpwYAAACoElEQVRIx61VPWhTYRQ9SZPm5yX9i41JXluTpqatGa4u6hAodVGxDqWIVrooKIitP50KgkJdUhCEDiJSEdsO7dKlDkVwsYIUH3YwSyjYJVAKjsUhy3GoTb73kjTPn+Eu37v3vPude+75ICKoFwMxzK4CNOAgffE5imh26lA/CS13NNDAfmwiWGSvZP4XuH7PXwY3oJG9MvhX4BfCGPsCV5b++AhF3IeBZ3QMrQFZ+vTb1ahSO9SutWP1YwnETUbSUzXBkxh85SyfsyG8QpGjVcGH2rFUBv5d0NY7VxO8G5Nr1vyG2ApFWk3gZxOYWLEkbiNYoC6DtWlB4mEQe4apzkG2pp4r4NDvB1BUk3bQlGfsQBGHDRT6gyAKau03NO+xe78WZ+LIqtfbhL/IThmxrZYk0tN+Cz3NqXkRAa60mv/MQHJpXyX2pXj6GGY/KBg/EM5RpBNTXrXQT8bVru3qHIlpjyWnR85DldM2QgWK9FgKtZvK4LbQtMseOWnNud6CXUOVcUd6CsvKdX4ismHVqohAujD0phG5HDw5hvrGq23jpTDyZXAnGUs/xaKjkis7q22N4TYV3EVG04/xzG2RkU1TsprbLZOcvWSXjOJGUD10ktH+mT8GT2F4WWHgYC44F8OagX+hBtpVEyUg3V3vKKJBUsi8cFmWQEvM230QBmJmjRvwkLET4yVvuRg2d2/AQTZG5qoqR+04jJcV5uXueH9gXqVFqTQhcAfePD1HJhlVtY+WJ02YeO0zb3aJ62R5wdShZGZ8ZgMzTJ7jKH52Vv9W8qRI/2jtl+g4Tj0KVHZUL3YQyFPvu2zjDYU2FkJ20V0fdBONRbraZxitPptDJXY3hOG3Xix8cmBj3YnCVzgLW2hY/w7vAl2hEfXVqRa/AGAchbvAF6o4AAAAAElFTkSuQmCC) left bottom no-repeat;
}

.icon-power.processing {
  background: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABcAAAAYCAYAAAARfGZ1AAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAAsTAAALEwEAmpwYAAACr0lEQVRIx61VXUhTYRj+GsloFYhUA0egUCZUfMuJCjoPYzEHWzO1WAbahRWSyBK8KWMSGNTFZEG1IUEKkrDIikmyixiWXRgkRjBDZxfDYNRseBFGBE/bWdv5ztnfSbp4b8553ud7v/d9nvcjlFJSLJqvQTcVgf3tV3RgZVgLSnfIySPFQSjpDsAejKQihDULOFr6v8iVF4Ms+TcLOmnZtshbPKh48h3HsewsT10/P3nDADSP1qFFzKdJYHcWIIfCehec/y/J4g/YcEd/OC/5aajdIeE7YuN1iQOUOclbbgrEmYTJVm1e8l5U+6T4mCd5QImYvA9VExJgFGtGDCavn68tUPXOoo3N4dU03XWMIYfycgAdLCiOsAHDaUUUGihUPS9gY3Mj+GxK/ye1V6Blr8cn36Bq2WrhUOqYl7THf66GJ7dOik/G+0s6sUmKS7HWAd1zhmMLL7nkcEn/PJu4Kalars6x1/FOgrlK9xF3mB3iayNMdLfUoXa/MJMIlk04SfdIMIrWKWG4vIyd+kPEy7YEj5tYKWWiG5rRRZhXAA6zFypyudHkglWkGnfTEfJA1KunnNQI8gKKUxMC+UIcbUkDkpElgTyMT2a5SylrufnZnv+0YIBqSBfTT/46Y+bqfyY/j4NekdZTcyGN12EIZsto13Zbws8u6qxPypnQRux3hSQm+DBUI/dB0A+jjjXhx9+wYqS5MrNbuFvi6vkD1u9pcyqHqdgyKibm8zbu16fXbxqoku6X1I7ZMiA6V4nbTJvaoeybRpVrTuzsTK8HBYMJlXSibCgoBks3XiCc/38Ym2Y8NJbnf4kSS6h/JruiYhHHugEe0wEZbygUZ8ag9S4VJ+WHtzF+FO25jVdQYmd9CdsvoOHZKkwzq7C9icK6FIfhC36dAF6p0VNo4JT8AfOI0HtZhtVRAAAAAElFTkSuQmCC) left bottom no-repeat;
}

.icon-download {
  margin: 0 8px;
  background: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAAsTAAALEwEAmpwYAAAA7klEQVQ4y2P4//8/AyUYhwQD6yIXhiAGDqkUBmG9TIbSG6GkGiBUzM/QySCsPB+MfQ/MBwqyDn8DGFh3ODPUwjXhwp13mnC7YA2Dtg0bks3o2H36fIbb/5XxewFoiCs/FkOcpsxmuPJfjrgwQDcEh2aIAasY5EQYGBwZRJS84dh1gdsRmCEgzasmyCWJoKkxzvVmOPNfkuFUOEMFVv9Grmtn2L9f5NxcBql6U4ZGrGoWvq1lOBSGL+S18MfIdIIGIGGl9PkJ3gwFDMLe8xnSnNy0GZi6wAbcKmFIZ+CVm02UISjYfD7D6k8pDJRmZwAKde5oo6iShAAAAABJRU5ErkJggg==) left bottom no-repeat;
}

.icon-pause {
  background: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAa0lEQVQ4T2NgoBT8/8/AutaDwY+BQyqFQVglE47NJqYAJfkJyYMMEKoVY+hkEFaej4qj5jM8/q9LSH7UgFEDhpEBvP0KWBRIpc9nuP9fg5A8JEfOZeAVV2YQYxCRl4Rgd0mGhv0i8ByLRx4AJGYUih69s0EAAAAASUVORK5CYII=) left bottom no-repeat;
}

.icon-resume {
  background: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAxklEQVQ4y2NgoAb4P4vBgEGpypth/38Rsgz4OI0hioFXbjaDUtR8hvyDUQz//7OSZ4Cw8nwwNsyezdB+wY18A2DYtX4iw7ynRuQbAMYG8xmcp9YyrPovR6YBMGw/nyF+aR7QIH4yDYBibAFNkgEwQ7b8dyTDAKA3AlYWY6QXwgZozWfwntjEsO2/MumBaFM2kWHKC3PSo9E8fT5D42VvolIligFKQfMZMnfGAzVykZYSWY2nMQQsKCYrQ/0vZOBkCMWSQIgAAA+ujguwc6ubAAAAAElFTkSuQmCC) left bottom no-repeat;
}

.icon-trash {
  background: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAACXBIWXMAAAsTAAALEwEAmpwYAAABUUlEQVQ4y2NgoAXYl8mgzsogZMEgqGTDIGhsw5CwVZdoze/WM3iKMghOZlAxnsYgotLDICLdwyBsPp9h4acUnJrmejCoAxUbgXC+FkMrg0HpfAYbDxuIGJdRuAJDJ4N903wGLQtzsJjLTF2G//9ZwZr//2fg7VdnqGXglZvNIKw1H6hgGk7MKjYNrMZ06kSgRkm4C4CGsC7TZWhkSDk1EW4yFnAqh6GCwWnDbKAaLhQJZAPO9TGYMFhPjMo5xyBlzuAezLBtr/JcbQY/hpn33Igy4NZ0hnQG7frZCVsZDMQZDPoZNh7ybpIEhkvf3faRaIACsgGdd5qIMqCwkIGTQXAVP0jMk0GIDxStOywZhBhWreLHa8AGU6ABQXunYUgiqdnhDUxwTjtmY1XzcCJDGA8oDwgrz8eNDeYzTHmYhzNP/J/LwMvAKyWCE4f+50fXAwB5Y94VTAfmBQAAAABJRU5ErkJggg==) left bottom no-repeat;
}

#footer {
  height: 30px;
  background-color: #2f2f2f;
  padding-top: 3px;
  text-align: center;
  font-size: 15px;
}

#footer a {
  color: #fdfdfd;
  font-weight: bold;
  letter-spacing: 1px;
}

@media (min-width: 400px) {
  .actions button.icon { margin: 0 20px 0 0; }
}

@media (min-width: 640px) {
  #header a.heading { float: none; }
  .actions button.icon { margin: 0 25px 0 0; }
}

@media (min-width: 960px) {
  .column.actions {
    order: 1;
    padding: 0 1rem 0 0;
  }

  .column.key {
    overflow: hidden;
    padding: 0 0 0 1rem;
  }

  .actions button.icon { margin: 0 0 0 30px; }
}

@media (min-width: 1122px) {
  #header { border-radius: 0 0 5px 5px; }

  #footer { border-radius: 5px 5px 0 0; }
}`

// index.html template string
var indexTmpl = `<!doctype html>
<html lang="en">
<head>
  <title>nogo</title>

  <!-- Metadata -->
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="theme-color" content="#9b4dca">

  <!-- Assets -->
  <link rel="stylesheet" href="//fonts.googleapis.com/css?family=Roboto:300,300italic,700,700italic">
  <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/normalize/5.0.0/normalize.min.css">
  <link rel="stylesheet" href="/css/nogo.css">
  <noscript>
    <style type="text/css">
      /* Power button and trash action requires JavaScript */
      #power-button, .actions button.icon-trash { display: none; }
    </style>
  </noscript>
</head>
<body>
  <header id="header" class="container">
    <div class="row">
      <div class="column">
        <a class="heading" href="/">nogo</a>
        <button id="power-button" class="icon icon-power"
          {{- if .isDisabled }}
          title="[Disabled] Click to enable nogo" data-disabled
          {{- else }}
          title="[Enabled] Click to disable nogo"
          {{- end -}}
        ></button>
      </div>
    </div>
  </header>

  <main id="main" class="container">
    <div id="inputs" class="row">
      <div class="column">
        <form action="/">
          <label for="search-input">Search Records</label>
          <input id="search-input" name="q" type="search" value="{{ .q }}" placeholder="Type part of a domain name, then press Enter." autocomplete="off" minlength="3" required>
        </form>
      </div>
      <div class="column">
        <form action="/records/" method="post">
          <label for="key-input">Add Record</label>
          <input id="key-input" name="key" type="text" value="" minlength="3" placeholder="Type a domain name, then press Enter." autocomplete="off" title="Must be a properly formatted domain." pattern=".+\..{2,}" required>
        </form>
      </div>
    </div>

    <div id="records-header" class="row">
      {{- if or .data .q .p }}
      <div id="back" class="column">
        <a href="/">&laquo; Back</a>
      </div>
      {{- else }}
      <div class="column">
        <a href="/?p=1">List Paused Records</a>
      </div>
      {{ end }}
      <div id="count" class="column text-right">
      {{- if or .data .q .p }}
        {{ if or .q }}Found {{ end }}<span id="data-count">{{ len .data }}</span> of <span id="total-count">{{ .totalCount }}</span> total records.
      {{- else }}
        <a class="icon icon-download" href="/export/hosts.txt" title="Download records as hosts file">&nbsp;</a>{{ .totalCount }} total records.
      {{- end }}
      </div>
    </div>

    {{- range $k, $v := .data }}
    <div id="{{ $k }}" class="row record">
      <div class="column actions"><!--
          Utilize a form for pause/resume, so that those with JavaScript
          disabled may be treated with equality.
     --><form action="/records/" method="post" class="pause-form {{ if $v.Paused }}hide{{ end }}">
          <input type="hidden" name="key" value="{{ $k }}" />
          <input type="hidden" name="paused" value="1" />
          <button class="icon icon-pause" type="submit" title="Pause" data-id="{{ $k }}"></button>
        </form><!-- clear white-space
     --><form action="/records/" method="post" class="resume-form {{ if not $v.Paused }}hide{{ end }}">
          <input type="hidden" name="key" value="{{ $k }}" />
          <input type="hidden" name="paused" value="0" />
          <button class="icon icon-resume" type="submit" title="Resume" data-id="{{ $k }}"></button>
        </form><!--
          JavaScript is required for delete, as there is no form method="delete"
          and I'm too much of a purist (or too apathetic) to resort to POST
          tunneling.
     --><button class="icon icon-trash" title="Delete" data-id="{{ $k }}"></button>
      </div>
      <div class="column key">{{ $k }}</div>
    </div>
    {{- end }}
  </main>

  <footer id="footer" class="container">
    <a href="http://nogo.curia.solutions/">http://nogo.curia.solutions</a>
  </footer>

  <script type="text/javascript">
    function togglePower() {
      var btn = document.getElementById('power-button');
      var disabled = btn.dataset.disabled === undefined ? true : false;

      btn.classList.add('processing');

      var req = new Request('/api/settings/', {
        method: 'PUT',
        body: JSON.stringify({ disabled: disabled })
      });

      fetch(req)
      .then(function(res) {
        btn.classList.remove('processing');

        if (res.ok) {
          // Toggle power button
          if (disabled) {
            btn.setAttribute('data-disabled', '');
            btn.title = '[Disabled] Click to enable nogo';
          } else {
            btn.removeAttribute('data-disabled');
            btn.title = '[Enabled] Click to disable nogo';
          }
        } else {
          // Shouldn't happen
          alert('ERROR: ' + res.status + ' ' + res.statusText);
        }
      });
    }

    function pauseRecord(key) {
      var req = new Request('/api/records/' + key, {
        method: 'PUT',
        body: JSON.stringify({paused: true})
      });

      fetch(req)
      .then(function(res) {
        if (res.ok) {
          // Hide pause, show resume
          document.getElementById(key).querySelector('.pause-form').classList.add('hide');
          document.getElementById(key).querySelector('.resume-form').classList.remove('hide');
        } else {
          // Shouldn't happen
          alert('ERROR: ' + res.status + ' ' + res.statusText);
        }
      });
    }

    function resumeRecord(key) {
      var req = new Request('/api/records/' + key, {
        method: 'PUT',
        body: JSON.stringify({paused: false})
      });

      fetch(req)
      .then(function(res) {
        if (res.ok) {
          // Hide resume, show pause
          document.getElementById(key).querySelector('.resume-form').classList.add('hide');
          document.getElementById(key).querySelector('.pause-form').classList.remove('hide');
        } else {
          // Shouldn't happen
          alert('ERROR: ' + res.status + ' ' + res.statusText);
        }
      });
    }

    function deleteRecord(key) {
      if (!confirm('Are you sure you want to delete this record?')) {
        return;
      }

      var req = new Request('/api/records/' + key, {method: 'DELETE'});

      fetch(req)
      .then(function(res) {
        if (res.ok) {
          // remove record
          document.getElementById(key).remove();

          // decrement counts
          document.getElementById('data-count').innerHTML = parseInt(document.getElementById('data-count').innerHTML) - 1;
          document.getElementById('total-count').innerHTML = parseInt(document.getElementById('total-count').innerHTML) - 1;
        } else {
          // Shouldn't happen
          alert('ERROR: ' + res.status + ' ' + res.statusText);
        }
      });
    }

    document.getElementById('power-button').addEventListener('click', function (evt) {
      togglePower();
      evt.preventDefault();
    });

    [].forEach.call(
      document.getElementsByClassName('icon-pause'),
      el => el.addEventListener('click', function (evt) {
        pauseRecord(this.dataset.id);
        evt.preventDefault();
      })
    );

    [].forEach.call(
      document.getElementsByClassName('icon-resume'),
      el => el.addEventListener('click', function (evt) {
        resumeRecord(this.dataset.id);
        evt.preventDefault();
      })
    );

    [].forEach.call(
      document.getElementsByClassName('icon-trash'),
      el => el.addEventListener('click', function () {
        deleteRecord(this.dataset.id);
      })
    );
  </script>
</body>
</html>`
