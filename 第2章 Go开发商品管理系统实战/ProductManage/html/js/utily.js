function RequestQueryString(key) {
   
     key = key.replace(/[\[]/, "\\\[").replace(/[\]]/, "\\\]");
     var regex = new RegExp("[\\?&]" + key + "=([^&#]*)");
     var qs = regex.exec(window.location.href);
     if (qs == null)
         return null;
     else
         return qs[1];
 }
