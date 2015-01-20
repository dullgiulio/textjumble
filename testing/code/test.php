<?php

class TestExampleClass {
	public function testMethod($arg1, $arg2) {
		echo "$arg1, $arg2"; // Just print the arguments
	}

	public function anotherTestMethod() {
		/*
		 * Define some array
		 */
		$arr = array(1, 2, 3, "string with many words");

		foreach ($arr as $a) {
			echo $a;
		
			return $a + $a;
		}
	}
}
