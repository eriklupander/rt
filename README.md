# rt
Go implementation of [The Ray Tracer challenge](https://pragprog.com/book/jbtracer/the-ray-tracer-challenge) book by Jamis Buck.

![image](renders/reference-highres-multisample-4x4-2depth-floor-refl.png)
_(image from a feature-branch with WiP multisampling and soft shadows)_

### Description
This is my WiP implementation of the Ray tracer as described in the book "The Ray Tracer Challenge" by https://pragprog.com/book/jbtracer/the-ray-tracer-challenge

### Changelog (performance fixes)
Just to keep track of which fixes were done when, kind of...

- 2020-02-20: Refactor to use arrays instead of slices for Matrices and Tuples.
- 2020-02-19: Use sort.Sort instead of sort.Slice
- 2020-02-15: Soft shadows and multisamling (feature-branch)
- 2020-02-05: Pass job per line instead of pixel
- 2020-01-10: Reduce allocations, cache stuff in render contexts
- 2019-12-28: Cache inverse
- 2019-12-25: Multi-threading

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
