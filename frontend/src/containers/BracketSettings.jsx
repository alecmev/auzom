import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';

import * as actions from '../actions';

import Settings from './Settings';

@connect(undefined, actions)
export default class BracketSettings extends Component {
  static propTypes = {
    stagePath: PropTypes.string,
    bracket: ImmutablePropTypes.map,

    updateBracket: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  static contextTypes = {
    router: PropTypes.object.isRequired,
  };

  handleSuccess = (m, x) => {
    m.slug && this.context.router.push(
      `${this.props.stagePath}/${x.slug}/_/settings`,
    );
    this.props.messagePush('bracket settings saved');
  };

  render() {
    return (
      <Settings
        isAdminOnly
        resourceName="bracket"
        resource={this.props.bracket}
        fields={{
          name: { name: 'Name' },
          abbr: { name: 'Abbreviation' },
          slug: { name: 'URL-friendly name' },
          order: { name: 'Sort number', type: 'number' },
          mapVetoProcedure: { name: 'Map veto procedure', type: 'textarea' },
        }}
        updateResource={this.props.updateBracket}
        onSuccess={this.handleSuccess}
      />
    );
  }
}
