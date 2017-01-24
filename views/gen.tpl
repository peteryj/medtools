<!DOCTYPE html>

<html>
<head>
  <title>MEDTools</title>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
</head>
<body>
  <header>
    <h4>Welcome to MED!!</h4>
    <div class="description">
      MED is a simple tool which is useful for writing papers. More features will be added!
    </div>
  </header>
  <form action="/" method="post">
    <p></p>
    <textarea id="papers" name="papers" rows="10" cols="100" placeholder="Please input your articles, one per line."></textarea>
    <br />
    <input type="submit" value="Get Citations" />
  </form>
  {{if .Result}}
  <table border="1">
    <th>Paper Title</th><th>Citations</th>
    {{range .Result}}
    <tr><td>{{.title}}</td><td>{{.cc}}</td></tr>
    {{end}}
  </table>
  {{end}}
  <footer>
    <div>
      Contact me:
      <a class="email" href="mailto:peter.yjzh@gmail.com">peter.yjzh@gmail.com</a>
    </div>
  </footer>
  <div class="backdrop"></div>
</body>
</html>
