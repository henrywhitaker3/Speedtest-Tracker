import { Id } from '../types';
export declare type State = Array<Id>;
export declare type Action = {
    type: 'ADD';
    toastId: Id;
    staleId?: Id;
} | {
    type: 'REMOVE';
    toastId?: Id;
};
export declare function reducer(state: State, action: Action): Id[];
