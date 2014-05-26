<?php
  // Create an image
require_once dirname(__FILE__).'/test.inc.php';

putenv('GDFONTPATH=' . realpath('.'));
$host = "Zabbix server";
$item = "system.cpu.load[percpu,avg1]";
$timeSpan = 3600;
if (!empty($_GET['host'])) {
  $host = $_GET['host'];
}
if (!empty($_GET['item'])) {
  $item = $_GET['item'];
}
if (!empty($_GET['span'])) {
  $timeSpan = $_GET['span'];
}
$sizeX = 600;
$sizeY = 250;
$dSizeX  = 500;
$dSizeY  = 150;
$dOrigoX = 75;
$dOrigoY = 175;

$maxTime  = time();
$minTime  = $maxTime - $timeSpan;
$vLineCnt = 11;
$minVal   = 0.1;
$maxVal   = 0.35;
$cntVal   = 11;
$hLineCnt = 6;

DBconnect();
$sql = "SELECT hi.clock AS clock, hi.value AS value".
       "  FROM items i, hosts h, history hi".
       "  WHERE h.hostid=i.hostid AND i.itemid=hi.itemid".
       "    AND h.host='".$host."'".
#       "    AND i.name='".$item."'".
       "    AND i.key_='".$item."'".
       "    AND hi.clock > Unix_Timestamp(now()) - ".$timeSpan.
       "  ORDER BY hi.clock";
$result = DBselect($sql);
$dbValues = array();
while ($dbValue = DBfetch($result)) {
  array_push($dbValues, $dbValue);
}
$sql = "SELECT min(hi.value) AS min, max(hi.value) AS max, count(*) AS cnt".
       "  FROM items i, hosts h, history hi".
       "  WHERE h.hostid=i.hostid AND i.itemid=hi.itemid".
       "    AND h.host='".$host."'".
#       "    AND i.name='".$item."'".
       "    AND i.key_='".$item."'".
       "    AND hi.clock > Unix_Timestamp(now()) - ".$timeSpan;
$result = DBfetch(DBselect($sql));
DBcommit();

$minVal = $result['min'];
$maxVal = $result['max'];
$cntVal = $result['cnt'];
$sql = "SELECT i.name AS name".
       "  FROM items i, hosts h".
       "  WHERE h.hostid=i.hostid".
       "    AND i.key_='".$item."'".
       "    AND h.host='".$host."'";
$result = DBfetch(DBselect($sql));
$varName = $result['name'];
DBcommit();

$header = $host.' : '.$varName;
$font   = 'DejaVuSans';
$image  = imagecreatetruecolor($sizeX, $sizeY);

// Allocate some colors
$white     = imagecolorallocate($image, 0xFF, 0xFF, 0xFF);
$gray      = imagecolorallocate($image, 0xC0, 0xC0, 0xC0);
$darkgray  = imagecolorallocate($image, 0x90, 0x90, 0x90);
$lightgray = imagecolorallocate($image, 0xE0, 0xE0, 0xE0);
$navy      = imagecolorallocate($image, 0x00, 0x00, 0x80);
$darknavy  = imagecolorallocate($image, 0x00, 0x00, 0x50);
$red       = imagecolorallocate($image, 0xFF, 0x00, 0x00);
$darkred   = imagecolorallocate($image, 0x90, 0x00, 0x00);
$black     = imagecolorallocate($image,    0,    0,    0);

imagestring($image, 4, $dOrigoX + 20, $dOrigoY - 20, $cntVal, $black);

// White background and a frame around all
imagefilledrectangle($image, 0, 0, $sizeX, $sizeY, $white);
imagerectangle($image, 0, 0, $sizeX - 1, $sizeY - 1, $navy);

// make header
imagestring($image, 4, 80, 1, $header, $black);

// Make the diagram area
$Xmax = $dOrigoX + $dSizeX;
$Ymin = $dOrigoY - $dSizeY;
imagefilledrectangle($image, $dOrigoX, $Ymin, $Xmax, $dOrigoY, $lightgray);
// Axis
imageline($image, $dOrigoX, $dOrigoY, $Xmax + 5, $dOrigoY, $black);
imageline($image, $dOrigoX, $dOrigoY, $dOrigoX, $Ymin - 5, $black);
// Arrows
imageline($image, $Xmax + 5, $dOrigoY -3 , $Xmax + 5,  $dOrigoY + 3, $black);
imageline($image, $Xmax + 5, $dOrigoY - 3, $Xmax + 10, $dOrigoY,     $black);
imageline($image, $Xmax + 5, $dOrigoY + 3, $Xmax + 10, $dOrigoY,     $black);
imageline($image, $dOrigoX - 3, $Ymin - 5, $dOrigoX + 3, $Ymin - 5,  $black);
imageline($image, $dOrigoX - 3, $Ymin - 5, $dOrigoX,     $Ymin - 10, $black);
imageline($image, $dOrigoX + 3, $Ymin - 5, $dOrigoX,     $Ymin - 10, $black);

// Draw time lines and write scale values
$Xspan  = $maxTime - $minTime;
$Xstep  = $Xspan / ($vLineCnt - 1);
$dStepX =  $dSizeX / ($vLineCnt - 1);
for ($i = 0; $i < $vLineCnt; $i++) {
  $Xpos = $dOrigoX + $i * $dStepX;
  if ($i > 0)
    imagedashedline($image, $Xpos, $dOrigoY, $Xpos, $Ymin, $darkgray);
  $str = date("H:i:s", $minTime + $i * $Xstep);
  imagestringup($image, 3, $Xpos - 7, $sizeY - 10, $str, $black);
}

// Make the horizontal lines
$Yspan  = $maxVal - $minVal;
$Ystep  = $Yspan / ($hLineCnt - 1);
$dStepY = $dSizeY / ($hLineCnt - 1);
for ($j = 0; $j < $hLineCnt; $j++) {
  $Ypos = $dOrigoY - $j * $dStepY;
  if ($j > 0)
    imageline($image, $dOrigoX, $Ypos, $Xmax, $Ypos, $darkgray);
  //  imagedashedline($image, $dOrigoX, $Ypos, $Xmax, $Ypos, $darkgray);
  $str = $minVal + $j * $Ystep;
  imagestring($image, 3, 10, $Ypos - 7, $str, $black);
}

// Draw the actual value
$scaleX = $dSizeX / $Xspan;
$scaleY = $dSizeY / $Yspan;

$prevX = -1;
$prevY = -1;
foreach ($dbValues as $dbValue) {
  //  $str = $dbValue['clock']." -> ".$dbValue['value'];
  //  print $str."\n";
  $timeVal = $dbValue['clock'];
  $measVal = $dbValue['value'];

  $plotX = (int) round($dOrigoX + ($timeVal - $minTime) * $scaleX);
  $plotY = (int) round($dOrigoY - ($measVal - $minVal) * $scaleY);

  if ($prevX > 0) {
    imageline($image, $prevX, $prevY, $plotX, $plotY, $red);
  }
  $prevX = $plotX;
  $prevY = $plotY;
}

  // Flush the image
header('Content-type: image/png');
imagepng($image);
imagedestroy($image);
?>
