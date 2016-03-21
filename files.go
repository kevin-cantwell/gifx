package palette

var PaletteHTML = `
<html>
<head>
  <style>
  body {
    margin: 0;
    padding: 0;
    font-family: 'courier new';
  }
  .content {
    clear: both;
  }
  .frame {
    border: 1px dashed grey;
    display: inline-block;
    position: relative;
  }
  .frame img {
    position: absolute;
  }
  .frame-meta {
    font-size: 14px;
    padding: 3px 0 3px 0;
    font-weight: bold;
    display: inline-block;
    vertical-align: top;
  }
  .rgba {
    color: white;
    text-shadow: -1px 0 black, 0 1px black, 1px 0 black, 0 -1px black;
    display: inline-block;
    width: 215px;
  }
  </style>
</head>
<body>
    {{range $i, $image := .}}
      {{if eq $i 0}}
        <div class="content">
          <div class="frame" style="width:{{$image.Orig.Config.Width}}px;height:{{$image.Orig.Config.Height}}px">
            <img src="{{$image.Image}}" />
          </div>
          <div class="frame-meta">
            <ul>
              <li>Frames: {{len $image.Orig.Image}}</li>
              <li>Loop Count: {{$image.Orig.LoopCount}}</li>
              <li>Width: {{$image.Orig.Config.Width}}</li>
              <li>Height: {{$image.Orig.Config.Height}}</li>
            </ul>
          </div>
          <br style="clear:both"/>
          <div class="palette"><!--
          {{range $color := $image.OrigPalette}}
            --><div class="rgba" style="background-color:rgba({{index $color.RGBA 0}},{{index $color.RGBA 1}},{{index $color.RGBA 2}},{{index $color.RGBA 3}})">
              rgba({{index $color.RGBA 0}},{{index $color.RGBA 1}},{{index $color.RGBA 2}},{{index $color.RGBA 3}})
            </div><!--
          {{end}}
          --></div>
        </div>
        <br style="clear:both"/>
        <br/>
        <br/>
      {{end}}
      <div class="content">
        <div class="frame" style="width:{{$image.Orig.Config.Width}}px;height:{{$image.Orig.Config.Height}}px">
          <img style="left:{{$image.X}};top:{{$image.Y}}" src="{{$image.Frame}}" />
        </div>
        <div class="frame-meta">
          <ul>
            <li>Frame Index: {{$i}}</li>
            <li>Delay: {{index $image.Orig.Delay $i}}</li>
            {{ $disposal := index $image.Orig.Disposal $i}}
            {{if eq $disposal 1}}
            <li>Disposal Method: None</li>
            {{end}}
            {{if eq $disposal 2}}
            <li>Disposal Method: None</li>
            {{end}}
            {{if eq $disposal 3}}
            <li>Disposal Method: Previous</li>
            {{end}}
            {{ $frame := index $image.Orig.Image $i}}
            <li>Bounds: {{print $frame.Bounds}}</li>
          </ul>
        </div>
        <br style="clear:both"/>
        <div class="palette"><!--
        {{range $color := $image.Palette}}
          --><div class="rgba" style="background-color:rgba({{index $color.RGBA 0}},{{index $color.RGBA 1}},{{index $color.RGBA 2}},{{index $color.RGBA 3}})">
            rgba({{index $color.RGBA 0}},{{index $color.RGBA 1}},{{index $color.RGBA 2}},{{index $color.RGBA 3}})
          </div><!--
        {{end}}
        --></div>
      </div>
      <br style="clear:both"/>
      <br/>
      <br/>
    {{end}}
</body>
</html>
`
