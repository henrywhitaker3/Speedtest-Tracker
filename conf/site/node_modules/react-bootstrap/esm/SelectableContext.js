import React from 'react';
var SelectableContext = React.createContext();
export var makeEventKey = function makeEventKey(eventKey, href) {
  if (eventKey != null) return String(eventKey);
  return href || null;
};
export default SelectableContext;