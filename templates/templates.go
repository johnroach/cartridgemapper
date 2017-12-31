package templates

//I haven't been able to find a better way to package templates within the binary
//if you find a better way please feel free to implement it.

// IndexPage template
var IndexPage = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-beta/css/bootstrap.min.css" integrity="sha384-/Y6pD6FV/Vv2HJnA6t+vslU6fwYXjCFtcEpHbNJ0lyAFsXTsjBbfaDjzALeQsN6M" crossorigin="anonymous">
    <style>
    /* Show it's not fixed to the top */
    body {
      min-height: 75rem;
    }
    </style>

  <link rel="stylesheet" type="text/css" href="https://cdn.datatables.net/1.10.16/css/jquery.dataTables.min.css">
  </head>
  <body>

    <nav class="navbar navbar-expand-md navbar-dark bg-dark mb-4">
          <a class="navbar-brand" href="#">Endeca Cartridge Mapper</a>
          <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarCollapse" aria-controls="navbarCollapse" aria-expanded="false" aria-label="Toggle navigation">
            <span class="navbar-toggler-icon"></span>
          </button>
          <div class="collapse navbar-collapse" id="navbarCollapse">
            <ul class="navbar-nav mr-auto">
              <li class="nav-item active">
                <a class="nav-link" href="#">Cartriges <span class="sr-only">
                (current)</span></a>
              </li>
              <li class="nav-item">
                <a class="nav-link" href="#">Sites</a>
              </li>
            </ul>
          </div>
        </nav>

        <div class="container">
          <table class="table" id="table_id">
            <thead class="thead-inverse">
              <tr>
                <th>Cartridge Name</th>
                <th>Cartridge Description</th>
                <th>Rules</th>
                <th>Sites</th>
                <th>Pages</th>
              </tr>
            </thead>
            <tbody>
              {{ range . }}
              <tr>
              <td>{{ .GetID }}</td>
              <td>{{ .GetDescription }}</td>
              <td>
                {{if .GetRules -}}
                  {{- range .GetRules }}
                    {{ . }}<br>
                  {{- end}}
                {{- else}}
                  No Rule found
                {{- end}}
              </td>
              <td>
                {{ if .GetSites -}}
                  {{- range .GetSites}}
                    {{ . }}<br>
                  {{- end}}
                {{- else}}
                  Cartridge not used in any site
                {{- end }}
              </td>
              <td>
              {{ if .GetPages -}}
                {{- range .GetPages}}
                  {{ . }}<br>
                {{- end}}
              {{- else}}
                Page is not used in any site
              {{- end }}
              </td>
              </tr>
              {{ end }}
            </tbody>
          </table>
        </div>

    <!-- Optional JavaScript -->
    <!-- jQuery first, then Popper.js, then Bootstrap JS -->
    <script src="https://code.jquery.com/jquery-3.2.1.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.11.0/umd/popper.min.js" integrity="sha384-b/U6ypiBEHpOf/4+1nzFpr53nxSS+GLCkfwBdFNTxtclqqenISfwAzpKaMNFNmj4" crossorigin="anonymous"></script>
    <script type="text/javascript" language="javascript" src="https://cdn.datatables.net/1.10.16/js/jquery.dataTables.min.js">
</script>
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-beta/js/bootstrap.min.js" integrity="sha384-h0AbiXch4ZDo7tp9hKZ4TsHbi047NrKGLO3SEJAg45jXxnGIfYzk4Si90RDIqNm1" crossorigin="anonymous"></script>
      <script>
    $(document).ready( function () {
      $('#table_id').DataTable();
    } );
    </script>
  </body>
</html>
`
