import { TestBed } from '@angular/core/testing';

import { ProviderAuthService } from './provider-auth.service';

describe("AuthService", () => {
  let service: ProviderAuthService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ProviderAuthService);
  });

  it("should be created", () => {
    expect(service).toBeTruthy();
  });
});
