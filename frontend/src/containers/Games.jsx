import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { Link } from 'react-router';

import * as actions from '../actions';
import * as selectors from '../selectors';
import Component from '../utils/Component';

@connect(createStructuredSelector({
  games: selectors.games,
}), actions)
export default class Games extends Component {
  static propTypes = {
    games: ImmutablePropTypes.map,
    loadGames: PropTypes.func.isRequired,
  };

  componentWillMount() {
    this.props.loadGames();
  }

  render() {
    const { games } = this.props;
    return (
      <div className={this.cni({ u: 'sectionSingle' })}>
        <div className={this.cn({ u: 'sectionMargined' })}>
          <ul>{games.toList().toJS().map(game => (
            <li key={game.id}>
              <Link to={`/${game.slug}`}>
                [{game.abbr}] {game.name}
              </Link>
            </li>
          ))}</ul>
        </div>
      </div>
    );
  }
}
