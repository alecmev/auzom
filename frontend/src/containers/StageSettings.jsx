import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';

import * as actions from '../actions';

import Settings from './Settings';

@connect(undefined, actions)
export default class StageSettings extends Component {
  static propTypes = {
    bracketsPath: PropTypes.string,
    stage: ImmutablePropTypes.map,

    updateStage: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  static contextTypes = {
    router: PropTypes.object.isRequired,
  };

  handleSuccess = (m, x) => {
    m.slug && this.context.router.push(
      `${this.props.bracketsPath}/${x.slug}/_/settings`,
    );
    this.props.messagePush('stage settings saved');
  };

  render() {
    return (
      <Settings
        isAdminOnly
        resourceName="stage"
        resource={this.props.stage}
        fields={{
          name: { name: 'Name' },
          abbr: { name: 'Abbreviation' },
          slug: { name: 'URL-friendly name' },
          startedAt: { name: 'Start time', type: 'datetime' },
        }}
        updateResource={this.props.updateStage}
        onSuccess={this.handleSuccess}
      />
    );
  }
}
