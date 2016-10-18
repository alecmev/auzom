import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';

import New from './New';

@connect(createStructuredSelector({
  gamesOptions: selectors.gamesOptions,
}), actions)
export default class TournamentNew extends Component {
  static propTypes = {
    gamesOptions: ImmutablePropTypes.list,

    loadGames: PropTypes.func.isRequired,
    createTournament: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  componentWillMount() {
    this.props.loadGames();
  }

  shouldComponentUpdate(nextProps) {
    return nextProps.gamesOptions !== this.props.gamesOptions;
  }

  handleSuccess = () => {
    // TODO: browserHistory.push(`/tournaments/${x.id}`);
    this.props.messagePush('created a tournament');
  };

  render() {
    return (
      <New
        isAdminOnly
        resourceName="tournament"
        fields={{
          gameId: {
            name: 'The game',
            type: 'select',
            options: this.props.gamesOptions.toJS(),
          },
          name: { name: 'Name' },
          abbr: { name: 'Abbreviation' },
          slug: { name: 'URL-friendly name' },
          blur: { name: 'Blur image' },
          logo: { name: 'Square logo image' },
          logoHasText: { name: 'Logo needs no tag overlay', type: 'boolean' },
          foundedAt: { type: 'date', name: 'Founding date' },
          description: { name: 'One-paragraph description', type: 'textarea' },
          email: { name: 'Support email', type: 'email' },
          twitch: { name: 'Twitch channel / team URL part' },
          youtube: { name: 'YouTube channel URL name' },
          twitter: { name: 'Twitter handle, without @' },
          facebook: { name: 'Facebook page URL name' },
          discord: { name: 'Discord invite code, not full link' },
          web: { name: 'Custom website link' },
          twitchLive: { name: 'LIVE Twitch channel URL name' },
        }}
        createResource={this.props.createTournament}
        onSuccess={this.handleSuccess}
      />
    );
  }
}
