import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';

import * as actions from '../actions';

import Editor from './Editor';

@connect(undefined, actions)
export default class SeasonRulesEditor extends Component {
  static propTypes = {
    season: ImmutablePropTypes.map,
    seasonPath: PropTypes.string,
    updateSeason: PropTypes.func.isRequired,
  };

  static contextTypes = {
    router: PropTypes.object.isRequired,
  };

  handleSuccess = () => this.context.router.push(
    `${this.props.seasonPath}/rules`,
  );

  render() {
    return (
      <Editor
        resourceName="season rules"
        resource={this.props.season}
        fields={{
          rules: {
            name: 'Rules',
            type: 'textarea',
            preview: 'text',
          },
        }}
        updateResource={this.props.updateSeason}
        onSuccess={this.handleSuccess}
      />
    );
  }
}
