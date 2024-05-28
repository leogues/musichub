import { DurationInMinutesPipe } from './duration-in-minutes.pipe';

describe('DurationInMinutesPipe', () => {
  it('create an instance', () => {
    const pipe = new DurationInMinutesPipe();
    expect(pipe).toBeTruthy();
  });
});
