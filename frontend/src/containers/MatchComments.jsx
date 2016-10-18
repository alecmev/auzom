import { PropTypes } from 'react';

import Component from '../utils/Component';

import Comments from './Comments';

export default class MatchComments extends Component {
  static propTypes = {
    matchId: PropTypes.string,
  };

  render() {
    return (
      <div className={this.cni({ u: 'sectionSingle' })}>
        <Comments
          target="match"
          targetId={this.props.matchId}
        />
      </div>
    );
  }
}
