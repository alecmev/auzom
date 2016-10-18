import { PropTypes } from 'react';

import * as utils from '../utils';

const Img = props => <img {...props} src={utils.https(props.src)} />;

Img.propTypes = {
  src: PropTypes.string.isRequired,
};

export default Img;
