let c = function( x ) {
  return function( m ) {
    return m ? function( y ) {
      return x;
    } : function( y ) {
      return (x = y);
    }
  }
};
let o = c( 10 );
o( false )( 15 );
o( true )( 42 );
