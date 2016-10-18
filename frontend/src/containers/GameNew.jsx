import { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';

import * as actions from '../actions';

import New from './New';

@connect(undefined, actions)
export default class GameNew extends Component {
  static propTypes = {
    createGame: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  handleSuccess = (x) => {
    browserHistory.push(`/${x.slug}`);
    this.props.messagePush('created a game');
  };

  render() {
    return (
      <New
        isAdminOnly
        resourceName="game"
        fields={{
          name: { name: 'Name' },
          abbr: { name: 'Abbreviation' },
          slug: { name: 'URL-friendly name' },
          cover: { name: 'Cover art URL' },
          summary: { name: 'Summary', type: 'textarea' },
          releasedAt: { type: 'date', name: 'Release date' },
        }}
        createResource={this.props.createGame}
        onSuccess={this.handleSuccess}
      />
    );
  }
}
