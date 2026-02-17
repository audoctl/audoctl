package fiberserver

var htmlScalar string = `<!DOCTYPE html>
<html>
  <body>
    <script
      id="api-reference"
      data-url="/api/swagger/swagger.json"
      src="https://cdnjs.cloudflare.com/ajax/libs/scalar-api-reference/1.36.2/standalone.min.js">
    </script>
  </body>
</html>`

var htmlRedoc = `<!DOCTYPE html>
<html>
  <head>
    <title>ApexScouty API Docs (Redoc)</title>
    <meta charset="utf-8"/>

    <!-- Redoc CDN -->
    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
    <style>
      body {margin: 0; padding: 0; }
      redoc { display:block; height:100vh; }
    </style>
  </head>

  <body>
    <redoc id="redoc" style="display:block; height:100vh;"></redoc>

    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
    <script>
      document.addEventListener("DOMContentLoaded", function() {
        Redoc.init('/api/swagger/swagger.json', {}, document.getElementById('redoc'));
      });
    </script>
  </body>
  </html>`

var htmlRapidoc = `<!DOCTYPE html>
<html>
  <body>
    <rapi-doc spec-url="/api/swagger/swagger.json"></rapi-doc>
    <script src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
  </body>
</html>`
