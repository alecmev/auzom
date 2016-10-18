import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';

import * as actions from '../actions';

import New from './New';

@connect(undefined, actions)
export default class StageNew extends Component {
  static propTypes = {
    season: ImmutablePropTypes.map,
    bracketsPath: PropTypes.string,

    createStage: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  static contextTypes = {
    router: PropTypes.object.isRequired,
  };

  handleSuccess = (x) => {
    this.context.router.push(`${this.props.bracketsPath}/${x.slug}`);
    this.props.messagePush('created a stage');
  };

  render() {
    if (!this.props.season) return null;
    return (
      <New
        isAdminOnly
        resourceName="stage"
        fields={{
          name: { name: 'Name' },
          abbr: { name: 'Abbreviation' },
          slug: { name: 'URL-friendly name' },
          startedAt: { name: 'Start time', type: 'datetime' },
          seasonId: { type: 'hidden', value: this.props.season.get('id') },
        }}
        createResource={this.props.createStage}
        onSuccess={this.handleSuccess}
      />
    );
  }
}
