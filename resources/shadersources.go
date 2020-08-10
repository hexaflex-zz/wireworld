package resources

var panelSources = [3]string{
	`#version 330 core

	layout (location = 0) in vec2 vPos;
	layout (location = 1) in vec2 vUV;
	
	uniform mat4 mvp;
	out vec2 fUV;
	
	void main()
	{
		gl_Position = mvp * vec4(vPos, 0.0, 1.0);
		fUV = vUV;
	}`,
	``,
	`#version 330 core

	uniform sampler2D img;
	
	in  vec2 fUV;
	out vec4 fragColor;
	
	void main()
	{
		fragColor = texture(img, fUV);
	}`,
}

var cellSelectorCellsSources = [3]string{
	`#version 330 core
	
	uniform mat4 mvp;
	
	layout (location = 0) in ivec2 cell;
	
	void main()
	{
		gl_Position = mvp * vec4(cell.x, cell.y, 0.0, 1.0);
	}`,
	`#version 330 core
	
	uniform vec2 cellSize;
	
	layout (points) in;
	layout (triangle_strip, max_vertices = 4) out;
	
	void main()
	{
		float cw = cellSize.x;
		float ch = cellSize.y;
	
		vec4 v = gl_in[0].gl_Position;
		vec4 a = vec4(v.x,    v.y,    0.0, 1.0);
		vec4 b = vec4(v.x,    v.y-ch, 0.0, 1.0);
		vec4 c = vec4(v.x+cw, v.y,    0.0, 1.0);
		vec4 d = vec4(v.x+cw, v.y-ch, 0.0, 1.0);
	
		gl_Position = a;
		EmitVertex();
		
		gl_Position = c;
		EmitVertex();
		
		gl_Position = b;
		EmitVertex();
		
		gl_Position = d;
		EmitVertex();
	
		EndPrimitive();	
	}`,
	`#version 330 core

	out vec4 fragColor;
	
	void main()
	{
		fragColor = vec4(1.0, 0.0, 0.596, 0.4);
	}`,
}

var cellSelectorRectSources = [3]string{
	`#version 330 core
	
	uniform mat4 mvp;
	
	layout (location = 0) in vec2 vPos;
	
	void main()
	{
		gl_Position = mvp * vec4(vPos, 0.0, 1.0);
	}`,
	``,
	`#version 330 core
	
	out vec4 fragColor;
	
	void main()
	{
		fragColor = vec4(0.0, 0.596, 1.0, 0.4);
	}`,
}

var gridSources = [3]string{
	`#version 330 core
	
	uniform mat4 mvp;
	
	layout (location = 0) in vec2 vPos;
	
	void main()
	{
		gl_Position = mvp * vec4(vPos, 0.0, 1.0);
	}`,
	``,
	`#version 330 core
	
	out vec4 fragColor;
	
	void main()
	{
		fragColor = vec4(0.75, 0.75, 0.75, 1.0);
	}`,
}

var cellRendererSources = [3]string{
	`#version 330 core
	
	uniform mat4 mvp;
	
	layout (location = 0) in ivec3 cell;
	
	flat out int gsColor;
	
	void main()
	{
		gl_Position = mvp * vec4(cell.x, cell.y, 0.0, 1.0);
		gsColor     = cell.z;
	}`,
	`#version 330 core
	
	uniform vec2 cellSize;
	
	flat in  int gsColor[];
	flat out int fsColor;
	
	layout (points) in;
	layout (triangle_strip, max_vertices = 4) out;
	
	void main()
	{
		float cw = cellSize.x;
		float ch = cellSize.y;
	
		vec4 v = gl_in[0].gl_Position;
		vec4 a = vec4(v.x,    v.y,    0.0, 1.0);
		vec4 b = vec4(v.x,    v.y-ch, 0.0, 1.0);
		vec4 c = vec4(v.x+cw, v.y,    0.0, 1.0);
		vec4 d = vec4(v.x+cw, v.y-ch, 0.0, 1.0);
	
		gl_Position = a;
		fsColor = gsColor[0];
		EmitVertex();
		
		gl_Position = c;
		EmitVertex();
		
		gl_Position = b;
		EmitVertex();
		
		gl_Position = d;
		EmitVertex();
	
		EndPrimitive();	
	}`,
	`#version 330 core
	
	uniform float alpha = 1.0;
	
	const vec4 palette[4] = vec4[](
		vec4(0.9,   0.9,   0.9, 1.0),
		vec4(1.0,   0.596, 0.0, 1.0),
		vec4(0.0,   0.596, 1.0, 1.0),
		vec4(0.596, 0.0,   1.0, 1.0)
	);
	
	flat in int fsColor;
	out vec4 fragColor;
	
	void main()
	{
		fragColor = palette[fsColor] * vec4(1, 1, 1, alpha);
	}`,
}
