import { createDevTools } from 'redux-devtools';
import Dispatcher from 'redux-devtools-dispatch';
import DockMonitor from 'redux-devtools-dock-monitor';
import LogMonitor from 'redux-devtools-log-monitor';
import MultipleMonitors from 'redux-devtools-multiple-monitors';

export default createDevTools(
  <DockMonitor
    toggleVisibilityKey="ctrl-q"
    changePositionKey="ctrl-h"
    fluid={false}
    defaultSize={300}
    defaultPosition="left"
    defaultIsVisible={false}
  >
    <MultipleMonitors>
      <LogMonitor theme="eighties" />
      <Dispatcher theme="eighties" />
    </MultipleMonitors>
  </DockMonitor>
);
