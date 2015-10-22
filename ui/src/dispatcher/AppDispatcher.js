/**
 * Created by ssoudan on 9/10/15.
 */

'use strict';

import {Dispatcher} from 'flux';
const instance: Dispatcher = new Dispatcher();
export default instance;

// So we can conveniently do, `import {dispatch} from './AppDispatcher';`
export const dispatch = instance.dispatch.bind(instance);