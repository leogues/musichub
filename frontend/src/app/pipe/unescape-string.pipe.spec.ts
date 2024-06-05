import { UnescapeStringPipe } from "./unescape-string.pipe";

describe('UnescapeStringPipe', () => {
  it('create an instance', () => {
    const pipe = new UnescapeStringPipe();
    expect(pipe).toBeTruthy();
  });
});
