package palette

var PaletteHTML = `
<html>
<head>
  <style>
  body {
    margin: 0;
    padding: 0;
    font-family: 'courier new';
    background: url('http://timehop.misc.s3.amazonaws.com/public/grid.gif');
  }
  .content {
    clear: both;
  }
  .frame {
    border: 1px dashed grey;
    vertical-align: top;
  }
  .frame-meta {
    font-size: 14px;
    padding: 3px 0 3px 0;
    font-weight: bold;
    display: inline-block;
  }
  .rgba {
    color: white;
    text-shadow: -1px 0 black, 0 1px black, 1px 0 black, 0 -1px black;
    float: left;
    width: 215px;
  }
  </style>
</head>
<body>
    {{range $i, $image := .}}
      {{if eq $i 0}}
        <div class="content">
          <img class="frame" src="{{$image.Image}}" />
          <div class="frame-meta">
            <ul>
              <li>Frames: {{len $image.Orig.Image}}</li>
              <li>Loop Count: {{$image.Orig.LoopCount}}</li>
              <li>Width: {{$image.Orig.Config.Width}}</li>
              <li>Height: {{$image.Orig.Config.Height}}</li>
            </ul>
          </div>
        </div>
        <br style="clear:both"/>
        <br/>
        <br/>
      {{end}}
      <div class="content">
        <img class="frame" src="{{$image.Frame}}" />
        <div class="frame-meta">
          <ul>
            <li>Index: {{$i}}</li>
            <li>Delay: {{index $image.Orig.Delay $i}}</li>
            {{ $disposal := index $image.Orig.Disposal $i}}
            {{if eq $disposal 1}}
            <li>Disposal: None</li>
            {{end}}
            {{if eq $disposal 2}}
            <li>Disposal: None</li>
            {{end}}
            {{if eq $disposal 3}}
            <li>Disposal: Previous</li>
            {{end}}
            {{ $frame := index $image.Orig.Image $i}}
            <li>Bounds: {{print $frame.Bounds}}</li>
          </ul>
        </div>
        <br style="clear:both"/>
        {{range $color := $image.Palette}}
          <span class="rgba" style="background-color:rgba({{index $color.RGBA 0}},{{index $color.RGBA 1}},{{index $color.RGBA 2}},{{index $color.RGBA 3}})">
            rgba({{index $color.RGBA 0}},{{index $color.RGBA 1}},{{index $color.RGBA 2}},{{index $color.RGBA 3}})
          </span>
        {{end}}
      </div>
      <br style="clear:both"/>
      <br/>
      <br/>
    {{end}}
</body>
</html>
`
