import { TestBed } from "@angular/core/testing";

import { MusicTableService } from "./music-table.service";

describe('MusicTableService', () => {
  let service: MusicTableService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(MusicTableService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
