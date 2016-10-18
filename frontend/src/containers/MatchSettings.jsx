import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';

import * as actions from '../actions';

import Settings from './Settings';

@connect(undefined, actions)
export default class MatchSettings extends Component {
  static propTypes = {
    match: ImmutablePropTypes.map,

    updateMatch: PropTypes.func.isRequired,
  };

  render() {
    return (
      <Settings
        lighter
        isAdminOnly
        resourceName="match"
        resource={this.props.match}
        fields={{
          startedAt: { name: 'Start time', type: 'datetime' },
          reportingClosedAt: { name: 'Report deadline', type: 'datetime' },
        }}
        updateResource={this.props.updateMatch}
      />
    );
  }
}
