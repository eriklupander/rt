# rt
Go implementation of [The Ray Tracer challenge](https://pragprog.com/book/jbtracer/the-ray-tracer-challenge) book by Jamis Buck.

![image](renders/reference-highres-multisample-4x4-2depth-floor-refl.png)
_(image from a feature-branch with WiP multisampling and soft shadows)_

### Description
This is my WiP implementation of the Ray tracer as described in the book "The Ray Tracer Challenge" by https://pragprog.com/book/jbtracer/the-ray-tracer-challenge

### Features, State etc
I think I've got all the features from the book covered, including:

- Spheres, Cubes, Cones, Cylinders, Planes, Triangles, CSGs
- Phong shading etc
- Multiple light sources
- Shadows, Reflection, Refraction, CSGs
- Groups
- Triangle render and .obj loading (possibly broken)

There's some unfinished / WiP stuff such as a YAML loader _NOT_ compatible with the .yaml format used in the book!

While all test cases from the books should be implemented, they may at times point at old (unused) implementations of certain features that I havn't gotten around to fix yet.

### Performance
Very naive implementation, has some basic caching of Inverse matrices, multi-threading using worker pool and some various measures to avoid allocating memory. Still loads of stuff to do since the allocation-heavy parts of the code really seems to kill performance.

### License 
MIT, see LICENSE.md
