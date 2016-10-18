import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';

import * as actions from '../actions';

import Settings from './Settings';

@connect(undefined, actions)
export default class TournamentSettings extends Component {
  static propTypes = {
    gameSlug: PropTypes.string.isRequired,
    tournament: ImmutablePropTypes.map,

    updateTournament: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  handleSuccess = (m, x) => {
    m.slug && browserHistory.push(`/${this.props.gameSlug}/${x.slug}/settings`);
    this.props.messagePush('tournament settings saved');
  };

  render() {
    return (
      <Settings
        lighter
        isAdminOnly
        resourceName="tournament"
        resource={this.props.tournament}
        fields={{
          name: { name: 'Name' },
          abbr: { name: 'Abbreviation' },
          slug: { name: 'URL-friendly name' },
          blur: { name: 'Blur image' },
          logo: { name: 'Square logo image' },
          logoHasText: { name: 'Logo has text in it', type: 'boolean' },
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
        updateResource={this.props.updateTournament}
        onSuccess={this.handleSuccess}
      />
    );
  }
}
