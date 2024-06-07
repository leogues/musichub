import { computed, Signal, signal, WritableSignal } from "@angular/core";

enum Status {
  pending = "pending",
  success = "success",
  error = "error",
}

export class APIQuery<T, E extends Error[] | Error | Record<string, Error>> {
  private defaultData: T;
  data: WritableSignal<T>;
  private _isFetching: WritableSignal<boolean>;
  private _status: WritableSignal<Status>;
  private _error: WritableSignal<E | null>;

  constructor(initialData: T) {
    this.defaultData = initialData;
    this.data = signal<T>(initialData);
    this._isFetching = signal<boolean>(false);
    this._status = signal<Status>(Status.pending);
    this._error = signal<E | null>(null);
  }

  get isFetching() {
    return this._isFetching.asReadonly();
  }

  get status() {
    return this._status.asReadonly();
  }

  get error() {
    return this._error.asReadonly();
  }

  private pending() {
    this._status.set(Status.pending);
  }

  fetching() {
    this._error.set(null);
    this._isFetching.set(true);
  }

  success(data: T) {
    this._status.set(Status.success);
    this.data.set(data);
    this._isFetching.set(false);
  }

  fail(error: E) {
    this._status.set(Status.error);
    this._error.set(error);
    this._isFetching.set(false);
  }

  restoreInitialState() {
    this.data.set(this.defaultData);
    this._isFetching.set(false);
    this.pending();
    this._error.set(null);
  }
}

export class DataQuery<T, E extends Error[] | Error | Record<string, Error>> {
  private APIQuery: APIQuery<T, E>;
  private _isError = computed(
    () => this.status() === Status.error || this.error() !== null,
  );
  private _isPending = computed(() => this.status() === Status.pending);
  private _isSuccessful = computed(() => this.status() === Status.success);

  constructor(APIQuery: APIQuery<T, E>) {
    this.APIQuery = APIQuery;
  }

  get data() {
    return this.APIQuery.data.asReadonly();
  }

  get status() {
    return this.APIQuery.status;
  }

  get error() {
    return this.APIQuery.error;
  }

  get isFetching() {
    return this.APIQuery.isFetching;
  }

  get isPending(): Signal<boolean> {
    return this._isPending;
  }

  get isError(): Signal<boolean> {
    return this._isError;
  }

  get isSuccessful(): Signal<boolean> {
    return this._isSuccessful;
  }
}
